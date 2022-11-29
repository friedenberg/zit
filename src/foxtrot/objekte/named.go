package objekte

import (
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
)

type Named[
	T gattung.Objekte,
	T1 gattung.ObjektePtr[T],
	T2 gattung.Identifier[T2],
	T3 gattung.IdentifierPtr[T2],
] struct {
	Stored  Stored[T, T1]
	Sha     sha.Sha
	Kennung T2
}

func (a *Named[T, T1, T2, T3]) Equals(b *Named[T, T1, T2, T3]) bool {
	if !a.Stored.Equals(&b.Stored) {
		return false
	}

	if !a.Kennung.Equals(&b.Kennung) {
		return false
	}

	return true
}

//TODO include kennung in reset
func (zn *Named[T, T1, T2, T3]) Reset() {
	// zn.Kennung = hinweis.Hinweis{}
	zn.Stored.Reset()
}
