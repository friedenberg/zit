package kasten

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type Objekte struct {
	Sha  sha.Sha
	Akte Akte
}

func (a Objekte) GetMetadatei() metadatei.Metadatei {
	return metadatei.Metadatei{
		AkteSha: a.Sha,
	}
}

func (a *Objekte) SetMetadatei(m metadatei.Metadatei) {
	a.Sha = m.AkteSha
}

func (a Objekte) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (o Objekte) Reset() {
	o.Sha = sha.Sha{}
	o.Akte.Reset()
}

func (o Objekte) ResetWith(o1 Objekte) {
	o.Sha = o1.Sha
	o.Akte.ResetWith(o1.Akte)
}

func (o Objekte) Equals(o1 Objekte) bool {
	if !o.Sha.Equals(o1.Sha) {
		return false
	}

	if !o.Akte.Equals(o1.Akte) {
		return false
	}

	return true
}

func (o Objekte) GetGattung() schnittstellen.Gattung {
	return gattung.Typ
}

func (o Objekte) GetAkteSha() schnittstellen.Sha {
	return o.Sha
}

func (o *Objekte) SetAkteSha(v schnittstellen.Sha) {
	o.Sha = sha.Make(v)
}
