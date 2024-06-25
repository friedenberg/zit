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
		&op.Kasten,
		op.Options,
		qg,
    func (col sku.CheckedOutLike) (err error) {
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

	return
}
