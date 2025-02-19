package inventory_list_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
)

func (store *Store) Commit(
	externalLike sku.ExternalLike,
	_ sku.CommitOptions,
) (err error) {
	sk := externalLike.GetSku()

	if sk.GetGenre() != genres.InventoryList {
		err = genres.MakeErrUnsupportedGenre(sk.GetGenre())
		return
	}

	if err = store.WriteInventoryListObject(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = store.ui.TransactedNew(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (store *Store) ReadOneInto(
	oid interfaces.ObjectId,
	_ *sku.Transacted,
) (err error) {
	if oid.GetGenre() != genres.InventoryList {
		err = genres.MakeErrUnsupportedGenre(oid.GetGenre())
		return
	}

	err = todo.Implement()
	// err = errors.BadRequestf("%q", oid)

	return
}

// TODO
func (s *Store) ReadPrimitiveQuery(
	queryGroup sku.PrimitiveQueryGroup,
	output interfaces.FuncIter[*sku.Transacted],
) (err error) {
	if err = s.ReadAllSkus(
		func(_, sk *sku.Transacted) (err error) {
			if err = output(sk); err != nil {
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
