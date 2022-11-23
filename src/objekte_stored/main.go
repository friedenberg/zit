package objekte_stored

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/sha"
)

type ObjektePtr[T any] interface {
	*T
	Equals(*T) bool
	Reset()
}

type Stored[T any, T1 ObjektePtr[T]] struct {
	Sha     sha.Sha
	Objekte T
	// Zettel zettel.Zettel
}

func (s *Stored[T, T1]) Reset() {
	s.Sha = sha.Sha{}
	T1(&(s.Objekte)).Reset()
}

func (a *Stored[T, T1]) Equals(b *Stored[T, T1]) bool {
	if !T1(&(a.Objekte)).Equals(&b.Objekte) {
		errors.Print("zettel")
		return false
	}

	if !a.Sha.Equals(b.Sha) {
		errors.Print("sha")
		return false
	}

	return true
}
