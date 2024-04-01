package remote_push

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
	"code.linenisgreat.com/zit/src/papa/remote_conn"
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

	var in []*sku.Transacted

	if err = d.Receive(&in); err != nil {
		err = errors.Wrap(err)
		return
	}

	out := make([]*sku.Transacted, 0)

	for _, sk := range in {
		// TODO-P2 support other Gattung
		if sk.GetGattung() != gattung.Zettel {
			continue
		}

		if op.umwelt.Standort().HasObjekte(
			op.umwelt.Konfig().GetStoreVersion(),
			sk.GetGattung(),
			sk.GetObjekteSha(),
		) {
			errors.Log().Printf("already have objekte: %s", sk.GetObjekteSha())
			return
		}

		errors.Log().Printf("need objekte: %s", sk.GetObjekteSha())

		// TODO-P1 check for akte sha
		// TODO-P1 write akte

		out = append(out, sk)
	}

	if err = d.Send(out); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
