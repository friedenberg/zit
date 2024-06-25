package user_ops

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

// TODO [radi/kof !task "add support for kasten in checkouts and external" project-2021-zit-features today zz-inbox]
type Edit struct {
	Kasten kennung.Kasten
	*umwelt.Umwelt
}

func (op Edit) Run(
	mode checkout_mode.Mode,
	zsc sku.CheckedOutLikeSet,
) (err error) {
	if err = op.GetStore().OpenFS(
		mode,
		op.PrinterHeader(),
		zsc,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = op.Reset(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var ms *query.Group

	builder := op.MakeQueryBuilderExcludingHidden(kennung.MakeGattung(gattung.Zettel))

	if ms, err = builder.WithCheckedOut(zsc).BuildQueryGroup(); err != nil {
		err = errors.Wrap(err)
		return
	}

	checkinOp := Checkin{}

	if err = checkinOp.Run(
		op.Umwelt,
		query.GroupWithKasten{
			Kasten: op.Kasten,
			Group:  ms,
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
