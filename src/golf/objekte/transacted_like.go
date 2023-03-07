package objekte

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
)

type TransactedLike interface {
	GetAkteSha() schnittstellen.Sha
	GetSku2() sku.Sku2
	GetSku() sku.Sku
	GetSkuLike() sku.SkuLike
	kennung.Matchable
	sku.DataIdentityGetter
}

type (
	FuncReaderTransacted[T TransactedLike] func(schnittstellen.FuncIter[T]) error
	FuncReaderTransactedLike               func(schnittstellen.FuncIter[TransactedLike]) error
)

type (
	FuncQuerierTransacted[T TransactedLike] func(kennung.Set, schnittstellen.FuncIter[T]) error
	FuncQuerierTransactedLike               func(kennung.Set, schnittstellen.FuncIter[TransactedLike]) error
)

func MakeApplyQueryTransactedLike[T TransactedLike](
	fat FuncQuerierTransacted[T],
) FuncQuerierTransactedLike {
	return func(ids kennung.Set, fatl schnittstellen.FuncIter[TransactedLike]) (err error) {
		return fat(
			ids,
			func(e T) (err error) {
				return fatl(e)
			},
		)
	}
}

func MakeApplyTransactedLike[T TransactedLike](
	fat FuncReaderTransacted[T],
) FuncReaderTransactedLike {
	return func(fatl schnittstellen.FuncIter[TransactedLike]) (err error) {
		return fat(
			func(e T) (err error) {
				return fatl(e)
			},
		)
	}
}
