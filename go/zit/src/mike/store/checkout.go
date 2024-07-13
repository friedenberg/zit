package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
)

func (s *Store) DeleteCheckout(col sku.CheckedOutLike) (err error) {
	kid := col.GetKasten().GetRepoIdString()
	es, ok := s.externalStores[kid]

	if !ok {
		err = errors.Errorf("no kasten with id %q", kid)
		return
	}

	if err = es.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = es.DeleteCheckout(col); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) CheckoutQuery(
	options checkout_options.Options,
	qg *query.Group,
	f interfaces.FuncIter[sku.CheckedOutLike],
) (err error) {
	qf := func(t *sku.Transacted) (err error) {
		var col sku.CheckedOutLike

		if col, err = s.CheckoutOne(qg.Kasten, options, t); err != nil {
			if errors.Is(err, sku.ErrExternalStoreUnsupportedTyp{}) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}

		sku.DetermineState(col, true)

		if err = s.checkedOutLogPrinter(col); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = f(col); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if err = s.QueryWithKasten(qg, qf); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) CheckoutOne(
	kasten kennung.Kasten,
	options checkout_options.Options,
	sz *sku.Transacted,
) (cz sku.CheckedOutLike, err error) {
	kid := kasten.GetRepoIdString()
	es, ok := s.externalStores[kid]

	if !ok {
		err = errors.Errorf("no kasten with id %q", kid)
		return
	}

	if cz, err = es.CheckoutOne(
		options,
		sz,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) UpdateCheckoutFromCheckedOut(
	options checkout_options.Options,
	col sku.CheckedOutLike,
) (err error) {
	switch col.GetKasten().GetRepoIdString() {
	case "chrome":
		err = todo.Implement()

	default:
		if err = s.cwdFiles.UpdateCheckoutFromCheckedOut(
			options,
			col,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
