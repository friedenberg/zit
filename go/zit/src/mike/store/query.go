package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
)

func (s *Store) QueryPrimitive(
	qg sku.PrimitiveQueryGroup,
	f interfaces.FuncIter[*sku.Transacted],
) (err error) {
	var e query.ExecutorPrimitive

	if e, err = s.MakeQueryExecutorPrimitive(qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = e.ExecuteTransacted(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

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

// TODO remove entirely in favor of Transacted
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

func (s *Store) ReadTransactedFromObjectId(
	k1 interfaces.ObjectId,
) (sk1 *sku.Transacted, err error) {
	sk1 = sku.GetTransactedPool().Get()

	if err = s.ReadOneInto(k1, sk1); err != nil {
		if collections.IsErrNotFound(err) {
			sku.GetTransactedPool().Put(sk1)
			sk1 = nil
		}

		err = errors.Wrap(err)
		return
	}

	return
}
