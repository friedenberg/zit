package remote_transfers

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
	"code.linenisgreat.com/zit/go/zit/src/papa/remote_conn"
)

type PullServer struct {
	env   *env.Local
	stage *remote_conn.StageSoldier
}

func MakePullServer(
	u *env.Local,
) (s PullServer, err error) {
	s = PullServer{
		env: u,
	}

	if s.stage, err = remote_conn.MakeStageSoldier(u); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.addToSoldierStage()

	return
}

func (op PullServer) Listen() (err error) {
	if err = op.stage.Listen(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (op PullServer) addToSoldierStage() {
	op.stage.RegisterHandler(
		remote_conn.DialogueTypeSkusForFilter,
		op.skusForFilter,
	)

	op.stage.RegisterHandler(
		remote_conn.DialogueTypeObjects,
		op.objectReaderForSku,
	)

	op.stage.RegisterHandler(
		remote_conn.DialogueTypeBlobs,
		op.blobReaderForSha,
	)
}

func (op PullServer) blobReaderForSha(
	d remote_conn.Dialogue,
) (err error) {
	defer errors.DeferredCloser(&err, d)

	var sh sha.Sha

	if err = d.Receive(&sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO-P2 rest is common

	ui.Log().Printf("received sha: %s", sh)

	var or io.ReadCloser

	if or, err = op.env.GetDirectoryLayout().BlobReader(&sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, or)

	var n int64

	if n, err = io.Copy(d, or); err != nil {
		err = errors.Wrap(err)
		return
	}

	ui.Log().Printf("served %d object bytes", n)

	return
}

func (op PullServer) objectReaderForSku(
	d remote_conn.Dialogue,
) (err error) {
	defer errors.DeferredCloser(&err, d)

	var msg messageRequestObjectData

	if err = d.Receive(&msg); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO-P2 rest is common

	ui.Log().Printf("received request: %#v", msg)

	orf := op.env.GetDirectoryLayout().ObjectReaderWriterFactory(msg.Gattung)

	var or io.ReadCloser

	if or, err = orf.ObjectReader(
		&msg.Sha,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, or)

	var n int64

	if n, err = io.Copy(d, or); err != nil {
		err = errors.Wrap(err)
		return
	}

	ui.Log().Printf("served %d object bytes", n)

	return
}

func (op PullServer) skusForFilter(
	d remote_conn.Dialogue,
) (err error) {
	defer errors.DeferredCloser(&err, d)

	var msg messageRequestSkus

	if err = d.Receive(&msg); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = op.env.GetStore().QueryTransacted(
		msg.MetaSet,
		quiter.MakeChain(
			func(sk *sku.Transacted) (err error) {
				if err = d.Send(sk); err != nil {
					err = errors.Wrap(err)
					return
				}

				return
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
