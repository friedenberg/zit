package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
)

func (s *Store) SaveBlob(el sku.ExternalLike) (err error) {
	repoId := el.GetRepoId()
	es, ok := s.externalStores[repoId]

	if !ok {
		err = errors.Errorf("no kasten with id %q", repoId)
		return
	}

	if err = es.SaveBlob(el); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

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

func (s *Store) Open(
	repoId ids.RepoId,
	m checkout_mode.Mode,
	ph interfaces.FuncIter[string],
	zsc sku.CheckedOutLikeSet,
) (err error) {
	es, ok := s.externalStores[repoId]

	if !ok {
		err = errors.Errorf("no repo id with id %q", repoId)
		return
	}

	if err = es.Open(m, ph, zsc); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) MakeQueryExecutor(
	qg *query.Group,
) (e query.Executor, err error) {
	if qg == nil {
		if qg, err = s.queryBuilder.BuildQueryGroup(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	es := s.externalStores[qg.RepoId]

	e = query.MakeExecutorWithExternalStore(
		qg,
		s.GetStreamIndex().ReadQuery,
		s.ReadOneInto,
		es,
	)

	return
}

func (s *Store) Merge(
	tm sku.Conflicted,
) (err error) {
	switch tm.CheckedOutLike.GetSkuExternalLike().GetRepoId().GetRepoIdString() {
	case "browser":
		err = todo.Implement()

	default:
		if err = s.cwdFiles.Merge(tm); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) RunMergeTool(
	tm sku.Conflicted,
) (err error) {
	tool := s.GetKonfig().ToolOptions.Merge

	switch tm.GetSkuExternalLike().GetRepoId().GetRepoIdString() {
	case "browser":
		err = todo.Implement()

	default:
		var co sku.CheckedOutLike

		if co, err = s.cwdFiles.RunMergeTool(tool, tm); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer s.PutCheckedOutLike(co)

		if err = s.CreateOrUpdateCheckedOut(co, false); err != nil {
			err = errors.Wrap(err)
			return
		}

	}

	return
}

func (s *Store) UpdateTransactedWithExternal(
	repoId ids.RepoId,
	z *sku.Transacted,
) (err error) {
	es, ok := s.externalStores[repoId]

	if !ok {
		err = errors.Errorf("no kasten with id %q", repoId)
		return
	}

	if err = es.UpdateTransacted(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadCheckedOutFromTransacted(
	kasten ids.RepoId,
	sk *sku.Transacted,
) (co sku.CheckedOutLike, err error) {
	switch kasten.GetRepoIdString() {
	case "browser":
		err = todo.Implement()

	default:
		if co, err = s.cwdFiles.ReadCheckedOutFromTransacted(sk); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
