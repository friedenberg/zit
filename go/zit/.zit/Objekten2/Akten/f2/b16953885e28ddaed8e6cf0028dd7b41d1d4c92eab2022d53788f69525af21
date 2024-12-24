package remote_push

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
	"code.linenisgreat.com/zit/go/zit/src/papa/remote_conn"
)

type Server struct {
	env   *env.Local
	stage *remote_conn.StageSoldier
}

func MakeServer(
	u *env.Local,
) (s Server, err error) {
	s = Server{
		env: u,
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
		remote_conn.DialogueTypeObjectWriter,
		op.ObjectWriter,
	)

	op.stage.RegisterHandler(
		remote_conn.DialogueTypeBlobWriter,
		op.BlobWriter,
	)
}

func (c Server) ObjectWriter(
	d remote_conn.Dialogue,
) (err error) {
	return
}

func (c Server) BlobWriter(
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
		if sk.GetGenre() != genres.Zettel {
			continue
		}

		if op.env.GetDirectoryLayout().HasObject(
			sk.GetGenre(),
			sk.GetObjectSha(),
		) {
			ui.Log().Printf("already have object: %s", sk.GetObjectSha())
			return
		}

		ui.Log().Printf("need object: %s", sk.GetObjectSha())

		// TODO-P1 check for blob sha
		// TODO-P1 write blob

		out = append(out, sk)
	}

	if err = d.Send(out); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
