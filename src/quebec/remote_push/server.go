package remote_push

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/golf/sku"
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
		remote_conn.DialogueTypeObjekteWriter,
		op.ObjekteWriter,
	)

	op.stage.RegisterHandler(
		remote_conn.DialogueTypeAkteWriter,
		op.AkteWriter,
	)
}

func (c Server) ObjekteWriter(
	d remote_conn.Dialogue,
) (err error) {
	return
}

func (c Server) AkteWriter(
	d remote_conn.Dialogue,
) (err error) {
	return
}

func (op Server) GetNeededSkus(
	d remote_conn.Dialogue,
) (err error) {
	defer errors.DeferredCloser(&err, d)

	var in []sku.Sku

	if err = d.Receive(&in); err != nil {
		err = errors.Wrap(err)
		return
	}

	out := make([]sku.Sku, 0)

	for _, sk := range in {
		//TODO-P2 support other Gattung
		if sk.Gattung != gattung.Zettel {
			continue
		}

		if op.umwelt.Standort().HasObjekte(sk.Gattung, sk.ObjekteSha) {
			errors.Log().Printf("already have objekte: %s", sk.ObjekteSha)
			return
		}

		errors.Log().Printf("need objekte: %s", sk.ObjekteSha)

		//TODO-P1 check for akte sha
		//TODO-P1 write akte

		out = append(out, sk)
	}

	if err = d.Send(out); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}