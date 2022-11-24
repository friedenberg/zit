package typ

import (
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/echo/konfig"
)

type Akte struct {
	KonfigTyp konfig.KonfigTyp
}

func (a *Akte) Reset(b *Akte) {
	if b == nil {
		a.KonfigTyp = konfig.KonfigTyp{}
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

func (a Akte) Gattung() gattung.Gattung {
	return gattung.Typ
}
