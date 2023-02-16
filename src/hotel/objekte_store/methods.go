package objekte_store

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/kennung"
)

func MethodForSigil[K any, T any](
	reader Querier[K, T],
	sigil kennung.Sigil,
) func(schnittstellen.FuncIter[T]) error {
	if sigil.IncludesHistory() {
		return reader.ReadAll
	} else {
		return reader.ReadAllSchwanzen
	}
}
