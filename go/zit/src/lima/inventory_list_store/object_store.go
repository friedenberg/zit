package inventory_list_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
)

func (s *Store) Commit(el sku.ExternalLike, _ sku.CommitOptions) (err error) {
	sk := el.GetSku()

	if sk.GetGenre() != genres.InventoryList {
		err = genres.MakeErrUnsupportedGenre(sk.GetGenre())
		return
	}

	if err = s.WriteInventoryListObject(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadOneInto(
	oid interfaces.ObjectId,
	_ *sku.Transacted,
) (err error) {
	if oid.GetGenre() != genres.InventoryList {
		err = genres.MakeErrUnsupportedGenre(oid.GetGenre())
		return
	}

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
