package query

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type (
	QueryCheckedOut interface {
		QueryCheckedOut(
			qg *Group,
			f interfaces.FuncIter[sku.CheckedOutLike],
		) (err error)
	}

	ExternalStore interface {
		sku.ExternalStoreUpdateTransacted
		QueryCheckedOut
	}

	ExecutionInfo struct {
		ExternalStore
		sku.FuncPrimitiveQuery
	}
)

type Executor struct {
	*Group
	ExecutionInfo
	Out interfaces.FuncIter[sku.ExternalLike]
}

func MakeExecutor(
	qg *Group,
	f sku.FuncPrimitiveQuery,
) Executor {
	return MakeExecutorWithExternalStore(qg, f, nil)
}

func MakeExecutorWithExternalStore(
	qg *Group,
	f sku.FuncPrimitiveQuery,
	es ExternalStore,
) Executor {
	return Executor{
		Group: qg,
		ExecutionInfo: ExecutionInfo{
			FuncPrimitiveQuery: f,
			ExternalStore:      es,
		},
	}
}

func (e *Executor) ExecuteCheckedOutLike(
	out interfaces.FuncIter[sku.CheckedOutLike],
) (err error) {
	if !e.dotOperatorActive {
		err = errors.Errorf("checked out queries must include dot operator")
		return
	}

	if err = e.executeExternalQueryCheckedOutLike(out); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (e *Executor) ExecuteExternalLike(
	out interfaces.FuncIter[sku.ExternalLike],
) (err error) {
	if e.dotOperatorActive {
		if err = e.executeExternalQueryExternalLike(out); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = e.executeInternalQueryExternalLike(out); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (e *Executor) ExecuteTransacted(
	out interfaces.FuncIter[*sku.Transacted],
) (err error) {
	out1 := func(el sku.ExternalLike) (err error) {
		sk := el.GetSku()

		if err = out(sk); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if e.dotOperatorActive {
		if err = e.executeExternalQueryExternalLike(out1); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = e.executeInternalQueryExternalLike(out1); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (e *Executor) executeExternalQueryCheckedOutLike(
	out interfaces.FuncIter[sku.CheckedOutLike],
) (err error) {
	if err = e.ExternalStore.QueryCheckedOut(
		e.Group,
		out,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (e *Executor) executeExternalQueryExternalLike(
	out interfaces.FuncIter[sku.ExternalLike],
) (err error) {
	if err = e.executeExternalQueryCheckedOutLike(
		func(col sku.CheckedOutLike) (err error) {
			z := col.GetSkuExternalLike()

			if err = out(z); err != nil {
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

func (e *Executor) executeInternalQueryExternalLike(
	out interfaces.FuncIter[sku.ExternalLike],
) (err error) {
	if err = e.FuncPrimitiveQuery(
		e.Group,
		e.makeEmitSkuSigilLatest(out),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (e *Executor) makeEmitSkuSigilLatest(
	out interfaces.FuncIter[sku.ExternalLike],
) interfaces.FuncIter[*sku.Transacted] {
	return func(z *sku.Transacted) (err error) {
		g := genres.Must(z.GetGenre())
		m, ok := e.Get(g)

		if !ok {
			return
		}

		if m.GetSigil().IncludesExternal() {
			if err = e.UpdateTransacted(z); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		if !m.ContainsSku(z) {
			return
		}

		if err = out(z); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}
