package remote_pull

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/golf/id_set"
	"github.com/friedenberg/zit/src/india/typ"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/oscar/umwelt"
	"github.com/friedenberg/zit/src/papa/remote_conn"
)

type Server struct {
	umwelt *umwelt.Umwelt
	stage  *remote_conn.StageSoldier
}

func MakeServer(
	u *umwelt.Umwelt,
) (s Server, err error) {
	s = Server{
		umwelt: u,
	}

	if s.stage, err = remote_conn.MakeStageSoldier(u); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.addToSoldierStage()

	return
}

func (op Server) Listen() (err error) {
	if err = op.stage.Listen(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (op Server) addToSoldierStage() {
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

func (op Server) akteReaderForSha(
	d remote_conn.Dialogue,
) (err error) {
	defer errors.DeferredCloser(&err, d)

	var sh sha.Sha

	if err = d.Receive(&sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Log().Printf("received sha: %s", sh)

	var or io.ReadCloser

	if or, err = op.umwelt.StoreObjekten().AkteReader(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, or.Close)

	var n int64

	if n, err = io.Copy(d, or); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Log().Printf("served %d objekte bytes", n)

	return
}

func (op Server) objekteReaderForSku(
	d remote_conn.Dialogue,
) (err error) {
	defer errors.DeferredCloser(&err, d)

	var msg messageRequestObjekteData

	if err = d.Receive(&msg); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Log().Printf("received request: %#v", msg)

	var or io.ReadCloser

	if or, err = op.umwelt.StoreObjekten().ObjekteReader(
		msg.Gattung,
		msg.Sha,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, or.Close)

	var n int64

	if n, err = io.Copy(d, or); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Log().Printf("served %d objekte bytes", n)

	return
}

func (op Server) skusForFilter(
	d remote_conn.Dialogue,
) (err error) {
	defer errors.DeferredCloser(&err, d)

	defer func() {
		if e := recover(); e != nil {
			errors.Todo(errors.P0, "panicked: %s", e)
			panic(e)
		}
	}()

	var filter id_set.Filter

	if err = d.Receive(&filter); err != nil {
		err = errors.Wrap(err)
		return
	}

	method := op.umwelt.StoreObjekten().Zettel().ReadAllSchwanzen

	if op.umwelt.Konfig().IncludeHistory {
		method = op.umwelt.StoreObjekten().Zettel().ReadAll
	}

	if err = op.umwelt.StoreObjekten().Typ().ReadAll(
		func(z *typ.Transacted) (err error) {
			sk := z.Sku.Sku2()

			if err = d.Send(sk); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = method(
		collections.MakeChain(
			zettel.WriterIds{Filter: filter}.WriteZettelVerzeichnisse,
			func(z *zettel.Transacted) (err error) {
				sk := z.Sku.Sku2()

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
