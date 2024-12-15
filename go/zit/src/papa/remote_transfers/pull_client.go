package remote_transfers

import (
	"net"
	"sync"
	"syscall"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
	"code.linenisgreat.com/zit/go/zit/src/papa/remote_conn"
)

type PullClient interface {
	SkusFromFilter(
		*query.Group,
		interfaces.FuncIter[*sku.Transacted],
	) error
	PullSkus(*query.Group) error
	Close() error
}

type client struct {
	env                *env.Local
	stage              *remote_conn.StageCommander
	chDone             chan struct{}
	chFilterSkuTickets chan struct{}
	common
}

func MakePullClient(u *env.Local, from string) (c *client, err error) {
	c = &client{
		env:                u,
		chDone:             make(chan struct{}),
		chFilterSkuTickets: make(chan struct{}, concurrentSkuFilterJobLimit),
		common: common{
			Local: u,
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
	ourVersion := u.GetConfig().GetImmutableConfig().GetStoreVersion()

	if ourVersion.Less(theirVersion) {
		err = errors.BadRequest(ErrPullRemoteHasHigherVersion)
		return
	}

	c.pmf = object_inventory_format.FormatForVersions(ourVersion, theirVersion)

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
	ids *query.Group,
	f interfaces.FuncIter[*sku.Transacted],
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

		ui.Log().Printf("received sku: %v", sk)

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
	f interfaces.FuncIter[*sku.Transacted],
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
		if quiter.IsStopIteration(err) {
			err = nil
		} else {
			ui.TodoP1("support net.ErrClosed downstream")
			err = errors.Wrap(err)
			errMulti.Add(err)
		}

		return
	}
}

func (c *client) ObjectReader(
	g interfaces.GenreGetter,
	sh interfaces.ShaGetter,
) (rc sha.ReadCloser, err error) {
	var d remote_conn.Dialogue

	if d, err = c.stage.StartDialogue(
		remote_conn.DialogueTypeObjects,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	msgRequest := messageRequestObjectData{
		Gattung: genres.Make(g.GetGenre()),
	}

	msgRequest.Sha.SetShaLike(sh)

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

func (c client) BlobReader(
	sh sha.ShaLike,
) (rc sha.ReadCloser, err error) {
	var d remote_conn.Dialogue

	if d, err = c.stage.StartDialogue(
		remote_conn.DialogueTypeBlobs,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = d.Send(sh.GetShaLike()); err != nil {
		err = errors.Wrap(err)
		return
	}

	var ow sha.WriteCloser

	if ow, err = c.env.GetDirectoryLayout().BlobWriter(); err != nil {
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
