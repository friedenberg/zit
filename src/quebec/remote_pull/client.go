package remote_pull

import (
	"net"
	"sync"
	"syscall"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/golf/id_set"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/oscar/umwelt"
	"github.com/friedenberg/zit/src/papa/remote_conn"
)

const (
	concurrentSkuFilterJobLimit = 100
	// concurrentSkuFilterJobLimit = 1
)

type FuncSku func(sku.Sku2) error

type Client interface {
	SkusFromFilter(id_set.Filter, FuncSku) error
	gattung.ObjekteReaderFactory
	gattung.AkteReaderFactory
	Close() error
}

type client struct {
	umwelt             *umwelt.Umwelt
	stage              *remote_conn.StageCommander
	chDone             chan struct{}
	chFilterSkuTickets chan struct{}
}

func MakeClient(u *umwelt.Umwelt, from string) (c *client, err error) {
	c = &client{
		chDone:             make(chan struct{}),
		chFilterSkuTickets: make(chan struct{}, concurrentSkuFilterJobLimit),
	}

	if c.stage, err = remote_conn.MakeStageCommander(
		u,
		from,
		"pull",
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c client) Close() (err error) {
	if err = c.stage.MainDialogue().Send(struct{}{}); err != nil {
		if errors.IsErrno(err, syscall.EPIPE) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	if err = c.stage.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c client) SkusFromFilter(ids id_set.Filter, f FuncSku) (err error) {
	var d remote_conn.Dialogue

	if d, err = c.stage.StartDialogue(
		remote_conn.DialogueTypeSkusForFilter,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	wg := &sync.WaitGroup{}
	errMulti := errors.MakeMulti()

	defer func() {
		d.Close()
		wg.Wait()

		errMulti.Add(err)

		if !errMulti.Empty() {
			err = errMulti
		}
	}()

	if err = d.Send(ids); err != nil {
		err = errors.Wrap(err)
		return
	}

	for {
		if !errMulti.Empty() {
			break
		}

		var sk sku.Sku2

		if err = d.Receive(&sk); err != nil {
			if errors.IsEOF(err) || errors.Is(err, net.ErrClosed) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}

		errors.Log().Printf("received sku: %v", sk)

		c.chFilterSkuTickets <- struct{}{}
		wg.Add(1)
		go c.makeAndProcessOneSkuWithFilter(
			sk,
			f,
			wg,
			errMulti,
		)
	}

	return
}

func (c *client) makeAndProcessOneSkuWithFilter(
	sk sku.Sku2,
	f FuncSku,
	wg *sync.WaitGroup,
	errMulti errors.Multi,
) {
	defer func() {
		<-c.chFilterSkuTickets
		if r := recover(); r != nil {
			//TODO-P0 add to err chan
			errors.Err().Printf("panicked during process one sku: %s", r)
		}

		wg.Done()
	}()

	if err := f(sk); err != nil {
		if collections.IsStopIteration(err) || errors.Is(err, net.ErrClosed) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		errMulti.Add(err)

		return
	}
}

func (c *client) ObjekteReader(
	g gattung.GattungLike,
	sh gattung.ShaLike,
) (rc sha.ReadCloser, err error) {
	var d remote_conn.Dialogue

	if d, err = c.stage.StartDialogue(
		remote_conn.DialogueTypeObjekten,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	msgRequest := messageRequestObjekteData{
		Gattung: g.GetGattung(),
		Sha:     sh.GetSha(),
	}

	if err = d.Send(msgRequest); err != nil {
		err = errors.Wrap(err)
		return
	}

	rc = sha.MakeReadCloser(d)

	return
}

func (c client) AkteReader(
	sh sha.Sha,
) (rc sha.ReadCloser, err error) {
	var d remote_conn.Dialogue

	if d, err = c.stage.StartDialogue(
		remote_conn.DialogueTypeAkten,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = d.Send(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	rc = sha.MakeReadCloser(d)

	return
}
