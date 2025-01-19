package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
)

func (s *Store) GetObjectStore() sku.ObjectStore {
	return s
}

func (s *Store) ReadPrimitiveQuery(
	qg sku.PrimitiveQueryGroup,
	w interfaces.FuncIter[*sku.Transacted],
) (err error) {
	return s.GetStreamIndex().ReadPrimitiveQuery(qg, w)
}
