package remote_transfers

import (
	"net"
	"sync"
	"syscall"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/oscar/umwelt"
	"github.com/friedenberg/zit/src/papa/remote_conn"
)

type PullClient interface {
	SkusFromFilter(
		matcher.Query,
		schnittstellen.FuncIter[*sku.Transacted],
	) error
	PullSkus(matcher.Query) error
	Close() error
}

type client struct {
	umwelt             *umwelt.Umwelt
	stage              *remote_conn.StageCommander
	chDone             chan struct{}
	chFilterSkuTickets chan struct{}
	common
}

func MakePullClient(u *umwelt.Umwelt, from string) (c *client, err error) {
	c = &client{
		umwelt:             u,
		chDone:             make(chan struct{}),
		chFilterSkuTickets: make(chan struct{}, concurrentSkuFilterJobLimit),
		common: common{
			Umwelt: u,
		},
	}

	if c.stage, err = remote_conn.MakeStageCommander(
		u,
		from,
		"pull",
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	theirVersion := c.stage.MainDialogue().GetAngeboren().GetStoreVersion()
	ourVersion := u.Konfig().GetAngeboren().GetStoreVersion()

	if ourVersion.Less(theirVersion) {
		err = errors.Normal(ErrPullRemoteHasHigherVersion)
		return
	}

	c.pmf = objekte_format.FormatForVersions(ourVersion, theirVersion)

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
	ids matcher.Query,
	f schnittstellen.FuncIter[*sku.Transacted],
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
		MetaSet: ids,
	}

	if err = d.Send(msg); err != nil {
		err = errors.Wrap(err)
		return
	}

	for {
		if !errMulti.Empty() {
			break
		}

		var sk *sku.Transacted

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
	sk *sku.Transacted,
	f schnittstellen.FuncIter[*sku.Transacted],
	wg *sync.WaitGroup,
	errMulti errors.Multi,
) {
	defer func() {
		// if r := recover(); r != nil {
		// 	errMulti.Add(errors.Errorf("panicked during process one sku: %s",
		// r))
		// }

		<-c.chFilterSkuTickets

		wg.Done()
	}()

	if err := f(sk); err != nil {
		if iter.IsStopIteration(err) {
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
		Sha:     sha.Make(sh.GetShaLike()),
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

	if err = d.Send(sh.GetShaLike()); err != nil {
		err = errors.Wrap(err)
		return
	}

	var ow sha.WriteCloser

	if ow, err = c.umwelt.Standort().AkteWriter(); err != nil {
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
