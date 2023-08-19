package objekte

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type (
	FuncReaderTransacted[T sku.SkuLike]       func(schnittstellen.FuncIter[T]) error
	FuncReaderTransactedPtr[T sku.SkuLikePtr] func(schnittstellen.FuncIter[T]) error
	FuncReaderTransactedLike                  func(schnittstellen.FuncIter[sku.SkuLike]) error
	FuncReaderTransactedLikePtr               func(schnittstellen.FuncIter[sku.SkuLikePtr]) error
)

type (
	FuncQuerierTransacted[T sku.SkuLike]       func(kennung.MatcherSigil, schnittstellen.FuncIter[T]) error
	FuncQuerierTransactedPtr[T sku.SkuLikePtr] func(kennung.MatcherSigil, schnittstellen.FuncIter[T]) error
	FuncQuerierTransactedLike                  func(kennung.MatcherSigil, schnittstellen.FuncIter[sku.SkuLike]) error
	FuncQuerierTransactedLikePtr               func(kennung.MatcherSigil, schnittstellen.FuncIter[sku.SkuLikePtr]) error
)

func MakeApplyQueryTransactedLikePtr[T sku.SkuLikePtr](
	fat FuncQuerierTransactedPtr[T],
) FuncQuerierTransactedLikePtr {
	return func(ids kennung.MatcherSigil, fatl schnittstellen.FuncIter[sku.SkuLikePtr]) (err error) {
		return fat(
			ids,
			func(e T) (err error) {
				return fatl(e)
			},
		)
	}
}

func MakeApplyTransactedLikePtr[T sku.SkuLikePtr](
	fat FuncReaderTransacted[T],
) FuncReaderTransactedLikePtr {
	return func(fatl schnittstellen.FuncIter[sku.SkuLikePtr]) (err error) {
		return fat(
			func(e T) (err error) {
				return fatl(e)
			},
		)
	}
}
