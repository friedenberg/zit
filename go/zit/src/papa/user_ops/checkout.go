package user_ops

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

type Checkout struct {
	checkout_options.Options
	*umwelt.Umwelt
}

func (op Checkout) Run(
	skus sku.TransactedSet,
) (zsc store_fs.CheckedOutMutableSet, err error) {
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
) (zsc store_fs.CheckedOutMutableSet, err error) {
	zsc = collections_value.MakeMutableValueSet[*store_fs.CheckedOut](nil)

	if err = op.Umwelt.GetStore().CheckoutQueryFS(
		op.Options,
		qg,
		iter.MakeAddClonePoolPtrFunc(
			zsc,
			store_fs.GetCheckedOutPool(),
			store_fs.CheckedOutResetter,
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
