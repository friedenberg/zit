package matcher

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type (
	FuncReaderTransactedLikePtr func(schnittstellen.FuncIter[*sku.Transacted]) error
	FuncSigilTransactedLikePtr  func(MatcherSigil, schnittstellen.FuncIter[*sku.Transacted]) error
	FuncQueryTransactedLikePtr  func(Query, schnittstellen.FuncIter[*sku.Transacted]) error
)

func MakeFuncReaderTransactedLikePtr(
	ms Query,
	fq FuncQueryTransactedLikePtr,
) FuncReaderTransactedLikePtr {
	return func(f schnittstellen.FuncIter[*sku.Transacted]) (err error) {
		return fq(ms, f)
	}
}