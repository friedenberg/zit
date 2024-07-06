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
	Kasten kennung.Kasten
	*umwelt.Umwelt
	checkout_options.Options
	Open    bool
	Edit    bool
	Utility string
}

func (op Checkout) Run(
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

	if zsc, err = op.RunQuery(qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (op Checkout) RunQuery(
	qg *query.Group,
) (zsc sku.CheckedOutLikeMutableSet, err error) {
	zsc = sku.MakeCheckedOutLikeMutableSet()
	var l sync.Mutex

	if err = op.Umwelt.GetStore().CheckoutQuery(
		op.Options,
		query.GroupWithKasten{
			Group:  qg,
			Kasten: op.Kasten,
		},
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
			op.Kasten,
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
				Kasten: op.Kasten,
				ExternalQuery: sku.ExternalQuery{
					Queryable: ms,
				},
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
