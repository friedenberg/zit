package store_util

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/kilo/objekte_store"
)

type GattungStoreLike interface{}

type CommonStoreBase struct {
	schnittstellen.GattungGetter
	StoreUtil
	objekte_store.TransactedReader
}

func MakeCommonStoreBase(
	gg schnittstellen.GattungGetter,
	sa StoreUtil,
	tr objekte_store.TransactedReader,
) (s *CommonStoreBase, err error) {
	s = &CommonStoreBase{
		GattungGetter:    gg,
		StoreUtil:        sa,
		TransactedReader: tr,
	}

	return
}

func (s *CommonStoreBase) Query(
	m matcher.MatcherSigil,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	return objekte_store.QueryMethodForMatcher(s, m, f)
}
