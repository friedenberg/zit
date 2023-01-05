package remote_pull

import (
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/golf/id_set"
	"github.com/friedenberg/zit/src/golf/sku"
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
		remote_conn.DialogueTypeAkteReaderForSha,
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

	var ar io.ReadCloser

	if ar, err = op.umwelt.StoreObjekten().AkteReader(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ar)

	var n int64

	if n, err = io.Copy(d, ar); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Log().Printf("served %d akte bytes", n)

	return
}

func (op Server) objekteReaderForSku(
	d remote_conn.Dialogue,
) (err error) {
	for {
		var strSku string

		if err = d.Receive(&strSku); err != nil {
			err = errors.Wrap(err)
			return
		}

		errors.Log().Printf("received sku str: %s", strSku)

		var sk sku.SkuLike

		if sk, err = sku.MakeSku(strSku); err != nil {
			err = errors.Wrap(err)
			return
		}

		errors.Log().Printf("received sku: %s", sk)

		var n int64

		if n, err = op.umwelt.StoreObjekten().SizeForObjektenSku(sk); err != nil {
			err = errors.Wrap(err)
			return
		}

		errors.Log().Printf("about to serve %d bytes", n)

		if _, err = io.WriteString(d, fmt.Sprintf("%d\n", n)); err != nil {
			err = errors.Wrap(err)
			return
		}

		func() {
			var or io.ReadCloser

			if or, err = op.umwelt.StoreObjekten().ReadCloserObjektenSku(sk); err != nil {
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
		}()
	}
}

func (op Server) skusForFilter(
	d remote_conn.Dialogue,
) (err error) {
	var filter id_set.Filter

	if err = d.Receive(&filter); err != nil {
		err = errors.Wrap(err)
		return
	}

	method := op.umwelt.StoreObjekten().Zettel().ReadAllSchwanzenVerzeichnisse

	if op.umwelt.Konfig().IncludeHistory {
		method = op.umwelt.StoreObjekten().Zettel().ReadAllVerzeichnisse
	}

	if err = method(
		collections.MakeChain(
			zettel.WriterIds{Filter: filter}.WriteZettelVerzeichnisse,
			func(z *zettel.Transacted) (err error) {
				if err = d.Send(sku.String(&z.Sku)); err != nil {
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

	if err = d.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
