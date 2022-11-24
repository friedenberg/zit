package typ

import (
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/konfig"
)

type Akte struct {
	Sha       sha.Sha
	KonfigTyp konfig.KonfigTyp
}

func (a *Akte) Reset(b *Akte) {
	if b == nil {
		a.Sha = sha.Sha{}
		a.KonfigTyp = konfig.KonfigTyp{}
	} else {
		a.Sha = b.Sha
		a.KonfigTyp = b.KonfigTyp
	}
}

func (a *Akte) Equals(b *Akte) bool {
	if !a.Sha.Equals(b.Sha) {
		return false
	}

	if !a.KonfigTyp.Equals(&b.KonfigTyp) {
		return false
	}

	return true
}

func (a Akte) Gattung() gattung.Gattung {
	return gattung.Typ
}
