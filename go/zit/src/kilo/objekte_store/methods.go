package objekte_store

import (
	"code.linenisgreat.com/zit-go/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit-go/src/bravo/iter"
	"code.linenisgreat.com/zit-go/src/hotel/sku"
	"code.linenisgreat.com/zit-go/src/india/matcher"
)

func QueryMethodForMatcher(
	reader TransactedReader,
	m matcher.MatcherSigil,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	out := reader.ReadAllSchwanzen

	if m.GetSigil().IncludesHistory() {
		out = reader.ReadAll
	}

	return out(
		iter.MakeChain(
			matcher.MakeMatcherFuncIter[*sku.Transacted](m),
			f,
		),
	)
}
