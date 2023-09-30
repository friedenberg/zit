package remote_transfers

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/oscar/umwelt"
	"github.com/friedenberg/zit/src/papa/remote_conn"
)

type PullServer struct {
	umwelt *umwelt.Umwelt
	stage  *remote_conn.StageSoldier
}

func MakePullServer(
	u *umwelt.Umwelt,
) (s PullServer, err error) {
	s = PullServer{
		umwelt: u,
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
		remote_conn.DialogueTypeObjekten,
		op.objekteReaderForSku,
	)

	op.stage.RegisterHandler(
		remote_conn.DialogueTypeAkten,
		op.akteReaderForSha,
	)
}

func (op PullServer) akteReaderForSha(
	d remote_conn.Dialogue,
) (err error) {
	defer errors.DeferredCloser(&err, d)

	var sh sha.Sha

	if err = d.Receive(&sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO-P2 rest is common

	errors.Log().Printf("received sha: %s", sh)

	var or io.ReadCloser

	if or, err = op.umwelt.Standort().AkteReader(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, or)

	var n int64

	if n, err = io.Copy(d, or); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Log().Printf("served %d objekte bytes", n)

	return
}

func (op PullServer) objekteReaderForSku(
	d remote_conn.Dialogue,
) (err error) {
	defer errors.DeferredCloser(&err, d)

	var msg messageRequestObjekteData

	if err = d.Receive(&msg); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO-P2 rest is common

	errors.Log().Printf("received request: %#v", msg)

	orf := op.umwelt.Standort().ObjekteReaderWriterFactory(msg.Gattung)

	var or io.ReadCloser

	if or, err = orf.ObjekteReader(
		msg.Sha,
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

	errors.Log().Printf("served %d objekte bytes", n)

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

	if err = op.umwelt.StoreObjekten().Query(
		msg.MetaSet,
		iter.MakeChain(
			// zettel.MakeWriterKonfig(
			// 	op.umwelt.Konfig(),
			// 	op.umwelt.StoreObjekten().GetAkten().GetTypV0(),
			// ),
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
