package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
)

func (s *Store) Query(
	qg sku.QueryGroup,
	f interfaces.FuncIter[*sku.Transacted],
) (err error) {
	if qg == nil {
		if qg, err = s.queryBuilder.BuildQueryGroup(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = s.GetVerzeichnisse().ReadQuery(
		qg,
		qg.MakeEmitSkuSigilLatest(
			f,
			ids.RepoId{},
			s.UpdateTransactedWithExternal,
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

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

	wg.Do(func() (err error) {
		if err = s.GetVerzeichnisse().ReadQuery(
			qg,
			qg.MakeEmitSkuMaybeExternal(
				f,
				qg.RepoId,
				s.UpdateTransactedWithExternal,
			),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	})

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
