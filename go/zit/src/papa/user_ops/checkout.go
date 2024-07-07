package user_ops

import (
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

type Checkout struct {
	*umwelt.Umwelt
	checkout_options.Options
	Open    bool
	Edit    bool
	Utility string
}

func (op Checkout) Run(
	skus sku.TransactedSet,
) (zsc sku.CheckedOutLikeMutableSet, err error) {
	var k kennung.Kasten

	if zsc, err = op.RunWithKasten(k, skus); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (op Checkout) RunWithKasten(
	kasten kennung.Kasten,
	skus sku.TransactedSet,
) (zsc sku.CheckedOutLikeMutableSet, err error) {
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

	if zsc, err = op.RunQuery(
		sku.ExternalQueryWithKasten{
			Kasten: kasten,
			ExternalQuery: sku.ExternalQuery{
				QueryGroup: qg,
			},
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (op Checkout) RunQuery(
	eqwk sku.ExternalQueryWithKasten,
) (zsc sku.CheckedOutLikeMutableSet, err error) {
	zsc = sku.MakeCheckedOutLikeMutableSet()
	var l sync.Mutex

	if err = op.Umwelt.GetStore().CheckoutQuery(
		op.Options,
		eqwk,
		func(col sku.CheckedOutLike) (err error) {
			l.Lock()
			defer l.Unlock()

			cl := col.Clone()

			if err = zsc.Add(cl); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if op.Utility != "" {
		eachAkteOp := EachAkte{
			Utility: op.Utility,
			Umwelt:  op.Umwelt,
		}

		if err = eachAkteOp.Run(zsc); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if op.Open || op.Edit {
		if err = op.GetStore().Open(
			eqwk.Kasten,
			op.CheckoutMode,
			op.PrinterHeader(),
			zsc,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if op.Edit {
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
			sku.ExternalQueryWithKasten{
				Kasten: eqwk.Kasten,
				ExternalQuery: sku.ExternalQuery{
					QueryGroup: ms,
				},
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
