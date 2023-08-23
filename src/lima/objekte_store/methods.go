package objekte_store

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/india/matcher"
)

type Query interface {
	schnittstellen.GattungGetter
	IncludesSchwanzen() bool
	IncludesHistory() bool
	IncludesCwd()
	// ContainsMatchable(kennung.Matchable) bool
}

func QueryMethodForMatcher[K any, T matcher.Matchable](
	reader Querier[K, T],
	m matcher.MatcherSigil,
	f schnittstellen.FuncIter[T],
) (err error) {
	out := reader.ReadAllSchwanzen

	if m.GetSigil().IncludesHistory() {
		out = reader.ReadAll
	}

	return out(
		iter.MakeChain(
			matcher.MakeMatcherFuncIter[T](m),
			f,
		),
	)
}
