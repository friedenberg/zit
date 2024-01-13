package objekte_store

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
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
