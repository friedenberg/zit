package query

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
)

type (
	QueryCheckedOut interface {
		QueryCheckedOut(
			query *Query,
			output interfaces.FuncIter[sku.SkuType],
		) (err error)
	}

	ExternalStore interface {
		interfaces.WorkspaceStoreReadAllExternalItems
		sku.ExternalStoreUpdateTransacted
		sku.ExternalStoreReadExternalLikeFromObjectIdLike
		QueryCheckedOut
	}

	WorkspaceEnv interface {
		IsTemporary() bool
		InWorkspace() bool
	}

	ExecutionInfo struct {
		ExternalStore
		sku.FuncPrimitiveQuery
		sku.FuncReadOneInto
	}
)

// TODO use ExecutorPrimitive
type Executor struct {
	primitive
	ExecutionInfo
	Out interfaces.FuncIter[sku.ExternalLike]
}

func MakeExecutorWithExternalStore(
	query *Query,
	fpq sku.FuncPrimitiveQuery,
	froi sku.FuncReadOneInto,
	externalStore ExternalStore,
	envWorkspace WorkspaceEnv,
) Executor {
	executor := Executor{
		primitive: primitive{query},
		ExecutionInfo: ExecutionInfo{
			FuncPrimitiveQuery: fpq,
			FuncReadOneInto:    froi,
		},
	}

	if envWorkspace.InWorkspace() {
		executor.ExternalStore = externalStore
	}

	return executor
}

// TODO refactor into methods that have internal in the name
func (executor *Executor) ExecuteExactlyOneExternalObject(
	permitInternal bool,
) (object *sku.Transacted, err error) {
	if executor.ExternalStore != nil {
		var externalObjectId ids.ObjectIdLike

		if externalObjectId, _, err = executor.Query.getExactlyOneExternalObjectId(
			permitInternal,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		object = sku.GetTransactedPool().Get()

		var external sku.ExternalLike

		if external, err = executor.ReadExternalLikeFromObjectIdLike(
			sku.CommitOptions{
				StoreOptions: sku.StoreOptions{
					UpdateTai: true,
				},
			},
			externalObjectId,
			object,
		); err != nil {
			err = errors.Wrapf(err, "ExternalObjectId: %q", externalObjectId)
			return
		}

		if external != nil {
			sku.TransactedResetter.ResetWith(object, external.GetSku())
		}
	} else {
		var objectId *ids.ObjectId

		if objectId, _, err = executor.Query.getExactlyOneObjectId(); err != nil {
			err = errors.Wrap(err)
			return
		}

		object = sku.GetTransactedPool().Get()

		if err = executor.FuncReadOneInto(
			objectId,
			object,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (executor *Executor) ExecuteExactlyOne() (sk *sku.Transacted, err error) {
	var objectId *ids.ObjectId
	var sigil ids.Sigil

	if objectId, sigil, err = executor.Query.getExactlyOneObjectId(); err != nil {
		err = errors.Wrap(err)
		return
	}

	sk = sku.GetTransactedPool().Get()

	if err = executor.ExecutionInfo.FuncReadOneInto(objectId, sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !sigil.IncludesExternal() {
		return
	}

	var external sku.ExternalLike

	if external, err = executor.ReadExternalLikeFromObjectIdLike(
		sku.CommitOptions{
			StoreOptions: sku.StoreOptions{
				UpdateTai: true,
			},
		},
		objectId,
		sk,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if external != nil {
		sku.TransactedResetter.ResetWith(sk, external.GetSku())
	}

	return
}

func (executor *Executor) ExecuteSkuType(
	out interfaces.FuncIter[sku.SkuType],
) (err error) {
	if executor.ExternalStore != nil {
		if err = executor.applyDotOperatorIfNecessary(); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = executor.executeExternalQueryCheckedOut(out); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = executor.executeInternalQuerySkuType(out); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (e *Executor) ExecuteTransacted(
	out interfaces.FuncIter[*sku.Transacted],
) (err error) {
	if err = e.readAllItemsIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO tease apart the reliance on dotOperatorActive here
	if e.dotOperatorActive && e.ExternalStore != nil {
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
	if err = e.readAllItemsIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if e.isDotOperatorActive() && e.ExternalStore != nil {
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
		e.Query,
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
		primitive{Query: e.Query},
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
		primitive{e.Query},
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
	return func(object *sku.Transacted) (err error) {
		g := genres.Must(object.GetGenre())

		if !e.containsSku(object) {
			return
		}

		// TODO cache query with sigil and object id
		genreQuery, ok := e.Get(g)

		if !ok {
			return
		}

		if genreQuery.GetSigil().IncludesExternal() && e.ExternalStore != nil {
			if err = e.UpdateTransacted(object); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		if err = out(object); err != nil {
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

		if !e.containsSku(internal) {
			return
		}

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

func (executor *Executor) applyDotOperatorIfNecessary() (err error) {
	if !executor.isDotOperatorActive() {
		return
	}

	if err = executor.readAllItemsIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (executor *Executor) readAllItemsIfNecessary() (err error) {
	if executor.ExternalStore == nil {
		return
	}

	if err = executor.ExternalStore.ReadAllExternalItems(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
