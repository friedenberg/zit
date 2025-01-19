package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
)

func (s *Store) GetObjectStore() external_store.ObjectStore {
	return s
}

func (s *Store) ReadPrimitiveQuery(
	qg sku.PrimitiveQueryGroup,
	w interfaces.FuncIter[*sku.Transacted],
) (err error) {
	return s.GetStreamIndex().ReadPrimitiveQuery(qg, w)
}
