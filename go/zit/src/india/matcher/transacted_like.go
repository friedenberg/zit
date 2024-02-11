package matcher

import (
	"code.linenisgreat.com/zit-go/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit-go/src/hotel/sku"
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
