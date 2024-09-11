package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
)

func (s *Store) QueryTransacted(
	qg *query.Group,
	f interfaces.FuncIter[*sku.Transacted],
) (err error) {
	var e query.Executor

	if e, err = s.MakeQueryExecutor(qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = e.ExecuteTransacted(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) Query(
	qg *query.Group,
	f interfaces.FuncIter[sku.ExternalLike],
) (err error) {
	var e query.Executor

	if e, err = s.MakeQueryExecutor(qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = e.ExecuteExternalLike(f); err != nil {
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

func (s *Store) QueryCheckedOut(
	qg *query.Group,
	f interfaces.FuncIter[sku.CheckedOutLike],
) (err error) {
	var e query.Executor

	if e, err = s.MakeQueryExecutor(qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = e.ExecuteCheckedOutLike(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
