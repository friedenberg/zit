package objekte_store

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/echo/kennung"
)

type Query interface {
	schnittstellen.GattungGetter
	IncludesSchwanzen() bool
	IncludesHistory() bool
	IncludesCwd()
	// ContainsMatchable(kennung.Matchable) bool
}

func QueryMethodForMatcher[K any, T kennung.Matchable](
	reader Querier[K, T],
	m kennung.MatcherSigil,
	f schnittstellen.FuncIter[T],
) (err error) {
	out := reader.ReadAllSchwanzen

	if m.GetSigil().IncludesHistory() {
		out = reader.ReadAll
	}

	return out(
		iter.MakeChain(
			kennung.MakeMatcherFuncIter[T](m),
			f,
		),
	)
}
