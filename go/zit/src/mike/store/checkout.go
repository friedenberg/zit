package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
)

func (s *Store) DeleteCheckedOutLike(col sku.CheckedOutLike) (err error) {
	if err = s.DeleteExternalLike(
		col.GetSkuExternalLike().GetRepoId(),
		col.GetSkuExternalLike(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) DeleteExternalLike(
	repoId ids.RepoId,
	el sku.ExternalLike,
) (err error) {
	es, ok := s.externalStores[repoId]

	if !ok {
		err = errors.Errorf("no kasten with id %q", repoId)
		return
	}

	if err = es.DeleteExternalLike(el); err != nil {
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
	es, ok := s.externalStores[qg.RepoId]

	if !ok {
		err = errors.Errorf("no kasten with id %q", qg.RepoId)
		return
	}

	qf := func(t *sku.Transacted) (err error) {
		var col sku.CheckedOutLike

    // TODO include a "query complete" signal for the external store to batch
    // the checkout if necessary
		if col, err = es.CheckoutOne(options, t); err != nil {
			if errors.Is(err, external_store.ErrUnsupportedTyp{}) {
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

	if err = s.QueryTransacted(qg, qf); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) CheckoutOne(
	repoId ids.RepoId,
	options checkout_options.Options,
	sz *sku.Transacted,
) (cz sku.CheckedOutLike, err error) {
	es, ok := s.externalStores[repoId]

	if !ok {
		err = errors.Errorf("no kasten with id %q", repoId)
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
	options checkout_options.OptionsWithoutMode,
	col sku.CheckedOutLike,
) (err error) {
	switch col.GetSkuExternalLike().GetRepoId().GetRepoIdString() {
	case "browser":
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
