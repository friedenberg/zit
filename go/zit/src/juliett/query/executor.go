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

	ExecutionInfo struct {
		sku.ExternalStoreUpdateTransacted
		sku.FuncPrimitiveQuery
		QueryCheckedOut
	}
)

type Executor struct {
	*Group
	ExecutionInfo
	Out interfaces.FuncIter[*sku.Transacted]
}

func (e *Executor) Execute() (err error) {
	if e.dotOperatorActive {
		if err = e.executeExternalQuery(); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = e.executeInternalQuery(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (e *Executor) executeExternalQuery() (err error) {
	if err = e.QueryCheckedOut.QueryCheckedOut(
		e.Group,
		func(col sku.CheckedOutLike) (err error) {
			z := col.GetSkuExternalLike().GetSku()

			if err = e.Out(z); err != nil {
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

func (e *Executor) executeInternalQuery() (err error) {
	if err = e.FuncPrimitiveQuery(
		e.Group,
		e.makeEmitSkuSigilLatest(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (qg *Group) makeEmitSku(
	f interfaces.FuncIter[*sku.Transacted],
) interfaces.FuncIter[*sku.Transacted] {
	return func(z *sku.Transacted) (err error) {
		g := genres.Must(z.GetGenre())
		m, ok := qg.Get(g)

		if !ok {
			return
		}

		if !m.ContainsSku(z) {
			return
		}

		if err = f(z); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func (e *Executor) makeEmitSkuSigilLatest() interfaces.FuncIter[*sku.Transacted] {
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

		if err = e.Out(z); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}
