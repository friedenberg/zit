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
