package inventory_list_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
)

func (s *Store) Commit(sku.ExternalLike, sku.CommitOptions) (err error) {
	err = todo.Implement()
	return
}

func (s *Store) ReadOneInto(interfaces.ObjectId, *sku.Transacted) (err error) {
	err = todo.Implement()
	return
}

func (s *Store) ReadPrimitiveQuery(
	qg sku.PrimitiveQueryGroup,
	w interfaces.FuncIter[*sku.Transacted],
) (err error) {
	err = todo.Implement()
	return
}
