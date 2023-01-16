package objekte

import (
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/golf/sku"
)

type TransactedLike interface {
	GetSku2() sku.Sku2
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
