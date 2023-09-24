package objekte

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
)

type (
	FuncReaderTransactedLikePtr func(schnittstellen.FuncIter[*sku.Transacted]) error
	FuncSigilTransactedLikePtr  func(matcher.MatcherSigil, schnittstellen.FuncIter[*sku.Transacted]) error
	FuncQueryTransactedLikePtr  func(matcher.Query, schnittstellen.FuncIter[*sku.Transacted]) error
)

func MakeFuncReaderTransactedLikePtr(
	ms matcher.Query,
	fq FuncQueryTransactedLikePtr,
) FuncReaderTransactedLikePtr {
	return func(f schnittstellen.FuncIter[*sku.Transacted]) (err error) {
		return fq(ms, f)
	}
}
