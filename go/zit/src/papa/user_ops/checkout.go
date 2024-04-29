package user_ops

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/juliett/query"
	"code.linenisgreat.com/zit/src/november/umwelt"
)

type Checkout struct {
	checkout_options.Options
	*umwelt.Umwelt
}

func (op Checkout) Run(
	skus sku.TransactedSet,
) (zsc sku.CheckedOutMutableSet, err error) {
	b := op.Umwelt.MakeQueryBuilder(
		kennung.MakeGattung(gattung.Zettel),
	).WithTransacted(
		skus,
	)

	var qg *query.Group

	if qg, err = b.BuildQueryGroup(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if zsc, err = op.RunQuery(qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (op Checkout) RunQuery(
	qg *query.Group,
) (zsc sku.CheckedOutMutableSet, err error) {
	zsc = collections_value.MakeMutableValueSet[*sku.CheckedOut](nil)
	add := func(co *sku.CheckedOut) (err error) {
		co1 := sku.GetCheckedOutPool().Get()
		sku.CheckedOutResetter.ResetWith(co1, co)

		return zsc.Add(co1)
	}

	if err = op.Umwelt.GetStore().CheckoutQuery(
		op.Options,
		qg,
		iter.MakeSyncSerializer(add),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
