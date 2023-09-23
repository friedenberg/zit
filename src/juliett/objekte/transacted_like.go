package objekte

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
)

type (
	FuncReaderTransacted[T sku.SkuLike]       func(schnittstellen.FuncIter[T]) error
	FuncReaderTransacted2                     func(schnittstellen.FuncIter[*sku.Transacted]) error
	FuncReaderTransactedPtr[T sku.SkuLikePtr] func(schnittstellen.FuncIter[T]) error
	FuncReaderTransactedLike                  func(schnittstellen.FuncIter[sku.SkuLike]) error
	FuncReaderTransactedLikePtr               func(schnittstellen.FuncIter[sku.SkuLikePtr]) error
)

type (
	FuncQuerierTransacted[T sku.SkuLike]       func(matcher.MatcherSigil, schnittstellen.FuncIter[T]) error
	FuncQuerierTransactedPtr[T sku.SkuLikePtr] func(matcher.MatcherSigil, schnittstellen.FuncIter[T]) error
	FuncQuerierTransactedPtr2                  func(matcher.MatcherSigil, schnittstellen.FuncIter[*sku.Transacted]) error
	FuncQuerierTransactedLike                  func(matcher.MatcherSigil, schnittstellen.FuncIter[sku.SkuLike]) error
	FuncQuerierTransactedLikePtr               func(matcher.MatcherSigil, schnittstellen.FuncIter[sku.SkuLikePtr]) error
)

func MakeApplyQueryTransactedLikePtr[T sku.SkuLikePtr](
	fat FuncQuerierTransactedPtr[T],
) FuncQuerierTransactedLikePtr {
	return func(ids matcher.MatcherSigil, fatl schnittstellen.FuncIter[sku.SkuLikePtr]) (err error) {
		return fat(
			ids,
			func(e T) (err error) {
				return fatl(e)
			},
		)
	}
}

func MakeApplyQueryTransactedLikePtr2[T sku.SkuLikePtr](
	fat FuncQuerierTransactedPtr2,
) FuncQuerierTransactedLikePtr {
	return func(ids matcher.MatcherSigil, fatl schnittstellen.FuncIter[sku.SkuLikePtr]) (err error) {
		return fat(
			ids,
			func(e *sku.Transacted) (err error) {
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

func MakeApplyTransactedLikePtr2(
	fat FuncReaderTransacted2,
) FuncReaderTransactedLikePtr {
	return func(fatl schnittstellen.FuncIter[sku.SkuLikePtr]) (err error) {
		return fat(
			func(e *sku.Transacted) (err error) {
				return fatl(e)
			},
		)
	}
}
