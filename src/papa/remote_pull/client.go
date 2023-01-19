package remote_pull

import (
	"net"
	"sync"
	"syscall"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/foxtrot/id_set"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/november/umwelt"
	"github.com/friedenberg/zit/src/oscar/remote_conn"
)

const (
	concurrentSkuFilterJobLimit = 100
)

type FuncSku func(sku.Sku2) error

type Client interface {
	SkusFromFilter(id_set.Filter, gattungen.Set, FuncSku) error
	PullSkus(id_set.Filter, gattungen.Set) error
	schnittstellen.ObjekteReaderFactory
	schnittstellen.AkteReaderFactory
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
		umwelt:             u,
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

func (c client) SkusFromFilter(
	ids id_set.Filter,
	gattungSet gattungen.Set,
	f FuncSku,
) (err error) {
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

	msg := messageRequestSkus{
		Filter:       ids,
		GattungSlice: gattungSet.Elements(),
	}

	if err = d.Send(msg); err != nil {
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
		if r := recover(); r != nil {
			errMulti.Add(errors.Errorf("panicked during process one sku: %s", r))
		}

		<-c.chFilterSkuTickets

		wg.Done()
	}()

	if err := f(sk); err != nil {
		if collections.IsStopIteration(err) {
			err = nil
		} else {
			errors.TodoP1("support net.ErrClosed downstream")
			err = errors.Wrap(err)
			errMulti.Add(err)
		}

		return
	}
}

func (c *client) ObjekteReader(
	g schnittstellen.GattungGetter,
	sh schnittstellen.ShaGetter,
) (rc sha.ReadCloser, err error) {
	var d remote_conn.Dialogue

	if d, err = c.stage.StartDialogue(
		remote_conn.DialogueTypeObjekten,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	msgRequest := messageRequestObjekteData{
		Gattung: gattung.Make(g.GetGattung()),
		Sha:     sha.Make(sh.GetSha()),
	}

	if err = d.Send(msgRequest); err != nil {
		if c.stage.ShouldIgnoreConnectionError(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	rc = sha.MakeReadCloser(d)

	return
}

func (c client) AkteReader(
	sh sha.ShaLike,
) (rc sha.ReadCloser, err error) {

	var d remote_conn.Dialogue

	if d, err = c.stage.StartDialogue(
		remote_conn.DialogueTypeAkten,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = d.Send(sh.GetSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	var ow sha.WriteCloser

	if ow, err = c.umwelt.StoreObjekten().AkteWriter(); err != nil {
		if c.stage.ShouldIgnoreConnectionError(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	rc = sha.MakeReadCloserTee(d, ow)

	return
}
