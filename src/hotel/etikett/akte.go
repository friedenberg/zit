package etikett

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
)

type Akte struct {
	Hide bool `toml:"hide"`
}

func (_ Akte) GetGattung() schnittstellen.GattungLike {
	return gattung.Etikett
}

func (ct *Akte) Merge(ct2 Akte) {
}

func (a *Akte) Reset() {
	a.Hide = false
}

func (a *Akte) ResetWith(b Akte) {
	a.Hide = b.Hide
}

func (a Akte) Equals(b Akte) bool {
	if a.Hide != b.Hide {
		return false
	}

	return true
}
