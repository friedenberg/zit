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
	var k kennung.RepoId

	if zsc, err = op.RunWithKasten(k, skus); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (op Checkout) RunWithKasten(
	kasten kennung.RepoId,
	skus sku.TransactedSet,
) (zsc sku.CheckedOutLikeMutableSet, err error) {
	b := op.Umwelt.MakeQueryBuilder(
		kennung.MakeGenre(gattung.Zettel),
	).WithTransacted(
		skus,
	)

	var qg *query.Group

	if qg, err = b.BuildQueryGroup(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if zsc, err = op.RunQuery(
		qg,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (op Checkout) RunQuery(
	eqwk *query.Group,
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
			eqwk.RepoId,
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

		builder := op.MakeQueryBuilderExcludingHidden(kennung.MakeGenre(gattung.Zettel))

		if ms, err = builder.WithCheckedOut(zsc).BuildQueryGroup(); err != nil {
			err = errors.Wrap(err)
			return
		}

		checkinOp := Checkin{}

		if err = checkinOp.Run(
			op.Umwelt,
			ms,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
