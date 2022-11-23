package objekte_stored

import (
	"github.com/friedenberg/zit/src/charlie/sha"
)

type Stored[T Objekte, T1 ObjektePtr[T]] struct {
	Sha     sha.Sha
	Objekte T
}

func (s *Stored[T, T1]) Reset() {
	s.Sha = sha.Sha{}
	T1(&(s.Objekte)).Reset(nil)
}

func (a *Stored[T, T1]) Equals(b *Stored[T, T1]) bool {
	if !T1(&(a.Objekte)).Equals(&b.Objekte) {
		return false
	}

	if !a.Sha.Equals(b.Sha) {
		return false
	}

	return true
}
