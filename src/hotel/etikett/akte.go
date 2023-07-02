package etikett

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
)

type Akte struct {
	AddToNewZettels bool `toml:"add-to-new-zettels"`
	Hide            bool `toml:"hide"`
}

func (_ Akte) GetGattung() schnittstellen.GattungLike {
	return gattung.Etikett
}

func (ct *Akte) Merge(ct2 Akte) {
}

func (a *Akte) Reset() {
	a.AddToNewZettels = false
	a.Hide = false
}

func (a *Akte) ResetWith(b Akte) {
	a.AddToNewZettels = b.AddToNewZettels
	a.Hide = b.Hide
}

func (a Akte) Equals(b Akte) bool {
	if a.AddToNewZettels != b.AddToNewZettels {
		return false
	}

	if a.Hide != b.Hide {
		return false
	}

	return true
}
