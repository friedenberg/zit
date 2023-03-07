package objekte_store

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type Query interface {
	schnittstellen.GattungGetter
	IncludesSchwanzen() bool
	IncludesHistory() bool
	IncludesCwd()
	// ContainsMatchable(kennung.Matchable) bool
}

func QueryMethodForSigil[K any, T any](
	reader Querier[K, T],
	sigil schnittstellen.IncludesHistory,
) func(schnittstellen.FuncIter[T]) error {
	if sigil.IncludesHistory() {
		return reader.ReadAll
	} else {
		return reader.ReadAllSchwanzen
	}
}
