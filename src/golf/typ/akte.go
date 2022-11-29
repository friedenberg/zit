package typ

import (
	"github.com/friedenberg/zit/src/delta/typ_toml"
)

type Akte struct {
	KonfigTyp typ_toml.Typ
}

func (a *Akte) Reset(b *Akte) {
	if b == nil {
		a.KonfigTyp = typ_toml.Typ{}
	} else {
		a.KonfigTyp = b.KonfigTyp
	}
}

func (a *Akte) Equals(b *Akte) bool {
	if !a.KonfigTyp.Equals(&b.KonfigTyp) {
		return false
	}

	return true
}
