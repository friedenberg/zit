package zettel

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type (
	MutableSet = schnittstellen.MutableSetLike[*sku.Transacted]
)

type TransactedUniqueKeyer struct{}

func (tk TransactedUniqueKeyer) GetKey(sz *sku.Transacted) string {
	if sz == nil {
		return ""
	}

	return collections.MakeKey(
		sz.GetTai(),
		sz.TransactionIndex,
		sz.GetKennung(),
		sz.ObjekteSha,
	)
}

type TransactedHinweisKeyer struct{}

func (tk TransactedHinweisKeyer) GetKey(sz *sku.Transacted) string {
	if sz == nil {
		return ""
	}

	return collections.MakeKey(
		sz.GetKennung(),
	)
}

func MakeMutableSetHinweis(c int) MutableSet {
	return collections_value.MakeMutableValueSet[*sku.Transacted](
		TransactedHinweisKeyer{},
	)
}
