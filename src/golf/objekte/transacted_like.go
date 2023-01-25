package objekte

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/foxtrot/sku"
)

type TransactedLike interface {
	GetAkteSha() schnittstellen.Sha
	GetSku2() sku.Sku2
	GetSku() sku.Sku
	GetSkuLike() sku.SkuLike
}

type FuncReaderTransacted[T TransactedLike] func(collections.WriterFunc[T]) error
type FuncReaderTransactedLike func(collections.WriterFunc[TransactedLike]) error

func MakeApplyTransactedLike[T TransactedLike](
	fat FuncReaderTransacted[T],
) FuncReaderTransactedLike {
	return func(fatl collections.WriterFunc[TransactedLike]) (err error) {
		return fat(
			func(e T) (err error) {
				return fatl(e)
			},
		)
	}
}
