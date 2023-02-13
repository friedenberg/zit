package collections

import "github.com/friedenberg/zit/src/alfa/schnittstellen"

type ValueSet2[E schnittstellen.Value] struct {
	schnittstellen.Set[E]
}

func (v ValueSet2[E]) String() string {
	return String[E](v.Set)
}
