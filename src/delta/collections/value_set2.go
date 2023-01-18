package collections

import "github.com/friedenberg/zit/src/schnittstellen"

type ValueSet2[E schnittstellen.Value] struct {
	SetLike[E]
}

func (v ValueSet2[E]) String() string {
	return String[E](v.SetLike)
}
