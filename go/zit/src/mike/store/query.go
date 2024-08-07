package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
)

func (s *Store) QueryWithKasten(
	qg *query.Group,
	f interfaces.FuncIter[*sku.Transacted],
) (err error) {
	if qg == nil {
		if qg, err = s.queryBuilder.BuildQueryGroup(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	wg := iter.MakeErrorWaitGroupParallel()
	es := s.externalStores[qg.RepoId.String()]

	e := &query.Executor{
		Group: qg,
		ExecutionInfo: query.ExecutionInfo{
			FuncPrimitiveQuery:            s.GetStreamIndex().ReadQuery,
			ExternalStoreUpdateTransacted: es,
			QueryCheckedOut:               es,
		},
		Out: func(el sku.ExternalLike) (err error) {
			if err = f(el.GetSku()); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	}

	wg.Do(e.Execute)

	if err = wg.GetError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) QueryWithKasten2(
	qg *query.Group,
	f interfaces.FuncIter[sku.ExternalLike],
) (err error) {
	if qg == nil {
		if qg, err = s.queryBuilder.BuildQueryGroup(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	wg := iter.MakeErrorWaitGroupParallel()
	es := s.externalStores[qg.RepoId.String()]

	e := &query.Executor{
		Group: qg,
		ExecutionInfo: query.ExecutionInfo{
			FuncPrimitiveQuery:            s.GetStreamIndex().ReadQuery,
			ExternalStoreUpdateTransacted: es,
			QueryCheckedOut:               es,
		},
		Out: f,
	}

	wg.Do(e.Execute)

	if err = wg.GetError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) QueryCheckedOut(
	qg *query.Group,
	f interfaces.FuncIter[sku.CheckedOutLike],
) (err error) {
	kid := qg.RepoId.GetRepoIdString()
	es, ok := s.externalStores[kid]

	if !ok {
		err = errors.Errorf("no kasten with id %q", kid)
		return
	}

	if err = es.QueryCheckedOut(
		qg,
		f,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
