package query

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/object_mode"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type (
	QueryCheckedOut interface {
		QueryCheckedOut(
			qg *Group,
			f interfaces.FuncIter[sku.SkuType],
		) (err error)
	}

	ExternalStore interface {
		sku.ExternalStoreReadAllExternalItems
		sku.ExternalStoreUpdateTransacted
		sku.ExternalStoreReadExternalLikeFromObjectId
		QueryCheckedOut
	}

	ExecutionInfo struct {
		ExternalStore
		sku.FuncPrimitiveQuery
		sku.FuncReadOneInto
	}
)

// TODO use ExecutorPrimitive
type Executor struct {
	*Group
	ExecutionInfo
	Out interfaces.FuncIter[sku.ExternalLike]
}

func MakeExecutor(
	qg *Group,
	fpq sku.FuncPrimitiveQuery,
	froi sku.FuncReadOneInto,
) Executor {
	return MakeExecutorWithExternalStore(qg, fpq, froi, nil)
}

func MakeExecutorWithExternalStore(
	qg *Group,
	fpq sku.FuncPrimitiveQuery,
	froi sku.FuncReadOneInto,
	es ExternalStore,
) Executor {
	return Executor{
		Group: qg,
		ExecutionInfo: ExecutionInfo{
			FuncPrimitiveQuery: fpq,
			FuncReadOneInto:    froi,
			ExternalStore:      es,
		},
	}
}

func (e *Executor) ExecuteExactlyOne() (sk *sku.Transacted, err error) {
	var k *ids.ObjectId
	var s ids.Sigil

	if k, s, err = e.Group.GetExactlyOneObjectId(
		genres.Zettel,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	sk = sku.GetTransactedPool().Get()

	if err = e.ExecutionInfo.FuncReadOneInto(k, sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !s.IncludesExternal() {
		return
	}

	var ze sku.ExternalLike

	if ze, err = e.ExecutionInfo.ReadExternalLikeFromObjectId(
		sku.CommitOptions{
			Mode: object_mode.ModeUpdateTai,
		},
		k,
		sk,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if ze != nil {
		sku.TransactedResetter.ResetWith(sk, ze.GetSku())
	}

	return
}

func (e *Executor) ExecuteSkuType(
	out interfaces.FuncIter[sku.SkuType],
) (err error) {
	if err = e.applyDotOperatorIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = e.executeExternalQueryCheckedOut(out); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (e *Executor) ExecuteTransacted(
	out interfaces.FuncIter[*sku.Transacted],
) (err error) {
	// TODO only apply dot operator when necessary
	if err = e.ExternalStore.ReadAllExternalItems(); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO tease apart the reliance on dotOperatorActive here
	if e.dotOperatorActive {
		if err = e.executeExternalQuery(out); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = e.executeInternalQuery(out); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (e *Executor) ExecuteTransactedAsSkuType(
	out interfaces.FuncIter[sku.SkuType],
) (err error) {
	// TODO only apply dot operator when necessary
	if err = e.ExternalStore.ReadAllExternalItems(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if e.DotOperatorActive() {
		if err = e.executeExternalQueryCheckedOut(out); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = e.executeInternalQuerySkuType(out); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (e *Executor) executeExternalQueryCheckedOut(
	out interfaces.FuncIter[sku.SkuType],
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

func (e *Executor) executeExternalQuery(
	out interfaces.FuncIter[*sku.Transacted],
) (err error) {
	if err = e.executeExternalQueryCheckedOut(
		func(col sku.SkuType) (err error) {
			z := col.GetSkuExternal()

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

func (e *Executor) executeInternalQuerySkuType(
	out interfaces.FuncIter[sku.SkuType],
) (err error) {
	if err = e.FuncPrimitiveQuery(
		e.Group,
		e.makeEmitSkuSigilLatestSkuType(out),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (e *Executor) executeInternalQuery(
	out interfaces.FuncIter[*sku.Transacted],
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
	out interfaces.FuncIter[*sku.Transacted],
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

func (e *Executor) makeEmitSkuSigilLatestSkuType(
	out interfaces.FuncIter[sku.SkuType],
) interfaces.FuncIter[*sku.Transacted] {
	return func(internal *sku.Transacted) (err error) {
		g := genres.Must(internal.GetGenre())
		m, ok := e.Get(g)

		if !ok {
			return
		}

		if m.GetSigil().IncludesExternal() {
			// TODO update External
			if err = e.UpdateTransacted(internal); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		if !m.ContainsSku(internal) {
			return
		}

		co := sku.GetCheckedOutPool().Get()
		defer sku.GetCheckedOutPool().Put(co)

		sku.TransactedResetter.ResetWith(co.GetSkuExternal(), internal)

		co.SetState(checked_out_state.Internal)

		if err = out(co); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func (e *Executor) applyDotOperatorIfNecessary() (err error) {
	if !e.DotOperatorActive() {
		return
	}

	if err = e.ExternalStore.ReadAllExternalItems(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
