package store_util

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/kilo/objekte_store"
)

type GattungStoreLike interface {
	Reindexer
}

type CommonStoreBase struct {
	schnittstellen.GattungGetter
	StoreUtil

	objekte_store.TransactedReader
	objekte_store.LogWriter
	persistentMetadateiFormat objekte_format.Format
}

func MakeCommonStoreBase(
	gg schnittstellen.GattungGetter,
	sa StoreUtil,
	tr objekte_store.TransactedReader,
	pmf objekte_format.Format,
) (s *CommonStoreBase, err error) {
	s = &CommonStoreBase{
		GattungGetter:             gg,
		StoreUtil:                 sa,
		TransactedReader:          tr,
		persistentMetadateiFormat: pmf,
	}

	return
}

func (s *CommonStoreBase) SetLogWriter(
	lw objekte_store.LogWriter,
) {
	s.LogWriter = lw
}

func (s *CommonStoreBase) Query(
	m matcher.MatcherSigil,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	return objekte_store.QueryMethodForMatcher(s, m, f)
}
