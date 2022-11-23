package objekte_stored

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/sha"
)

type Identifier[T any] interface {
	Equatable[T]
}

type Named[T any, T1 ObjektePtr[T], T2 Identifier[T2]] struct {
	Stored  Stored[T, T1]
	Sha     sha.Sha
	Kennung T2
}

func (a *Named[T, T1, T2]) Equals(b *Named[T, T1, T2]) bool {
	if !a.Stored.Equals(&b.Stored) {
		errors.Print("stored")
		return false
	}

	if !a.Kennung.Equals(&b.Kennung) {
		errors.Print("hinweis")
		return false
	}

	return true
}

func (zn *Named[T, T1, T2]) Reset() {
	// zn.Kennung = hinweis.Hinweis{}
	zn.Stored.Reset()
}
