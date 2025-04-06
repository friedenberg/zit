package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/env_workspace"
	pkg_query "code.linenisgreat.com/zit/go/zit/src/kilo/query"
)

func (s *Store) SaveBlob(el sku.ExternalLike) (err error) {
	es := s.envWorkspace.GetStore()

	if err = es.SaveBlob(el); err != nil {
		if errors.Is(err, env_workspace.ErrUnsupportedOperation{}) {
			err = nil
		} else {
			err = errors.Wrapf(err, "Sku: %s", el)
			return
		}
	}

	return
}

func (s *Store) DeleteCheckedOut(col *sku.CheckedOut) (err error) {
	es := s.envWorkspace.GetStore()

	if err = es.DeleteCheckedOut(col); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (store *Store) CheckoutQuery(
	options checkout_options.Options,
	query *pkg_query.Query,
	out interfaces.FuncIter[sku.SkuType],
) (err error) {
	externalStore := store.envWorkspace.GetStore()

	qf := func(t *sku.Transacted) (err error) {
		var co sku.SkuType

		// TODO include a "query complete" signal for the external store to batch
		// the checkout if necessary
		if co, err = externalStore.CheckoutOne(options, t); err != nil {
			if errors.Is(err, env_workspace.ErrUnsupportedType{}) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}

		if !store.envWorkspace.IsTemporary() {
			if err = store.ui.CheckedOutCheckedOut(co); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		if err = out(co); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if err = store.QueryTransacted(query, qf); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) CheckoutOne(
	repoId ids.RepoId,
	options checkout_options.Options,
	sz *sku.Transacted,
) (cz sku.SkuType, err error) {
	es := s.envWorkspace.GetStore()

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
	col sku.SkuType,
) (err error) {
	es := s.envWorkspace.GetStore()

	if err = es.UpdateCheckoutFromCheckedOut(
		options,
		col,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) Open(
	repoId ids.RepoId,
	m checkout_mode.Mode,
	ph interfaces.FuncIter[string],
	zsc sku.SkuTypeSet,
) (err error) {
	es := s.envWorkspace.GetStore()

	if err = es.Open(m, ph, zsc); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (store *Store) makeQueryExecutor(
	queryGroup *pkg_query.Query,
) (executor pkg_query.Executor, err error) {
	if queryGroup == nil {
		if queryGroup, err = store.queryBuilder.BuildQueryGroup(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	externalStore := store.envWorkspace.GetStore()

	executor = pkg_query.MakeExecutorWithExternalStore(
		queryGroup,
		store.GetStreamIndex().ReadPrimitiveQuery,
		store.ReadOneInto,
		externalStore,
	)

	return
}

// TODO make this configgable
func (store *Store) MergeConflicted(
	conflicted sku.Conflicted,
) (err error) {
	es := store.envWorkspace.GetStore()

	if err = es.Merge(conflicted); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) RunMergeTool(
	conflicted sku.Conflicted,
) (err error) {
	tool := s.GetConfig().GetCLIConfig().ToolOptions.Merge

	switch conflicted.GetSkuExternal().GetRepoId().GetRepoIdString() {
	case "browser":
		err = todo.Implement()

	default:
		var co sku.SkuType

		if co, err = s.envWorkspace.GetStoreFS().RunMergeTool(
			tool,
			conflicted,
		); err != nil {
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
	es := s.envWorkspace.GetStore()

	if err = es.UpdateTransacted(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadCheckedOutFromTransacted(
	repoId ids.RepoId,
	sk *sku.Transacted,
) (co *sku.CheckedOut, err error) {
	es := s.envWorkspace.GetStore()

	if co, err = es.ReadCheckedOutFromTransacted(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) UpdateTransactedFromBlobs(
	co *sku.CheckedOut,
) (err error) {
	external := co.GetSkuExternal()

	es := s.envWorkspace.GetStore()

	if err = es.UpdateTransactedFromBlobs(external); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
