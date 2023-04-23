package typ

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type Objekte struct {
	Akte
	Metadatei metadatei.Metadatei `toml:"-"`
}

func (a Objekte) GetMetadatei() metadatei.Metadatei {
	return a.Metadatei
}

func (a *Objekte) SetMetadatei(m metadatei.Metadatei) {
	a.Metadatei = m
}

func (a Objekte) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (o Objekte) Reset() {
	o.Akte.Reset()
	o.Metadatei.Reset()
}

func (o Objekte) ResetWith(o1 Objekte) {
	o.Metadatei.ResetWith(o1.Metadatei)
	o.Akte.ResetWith(o1.Akte)
}

func (o Objekte) Equals(o1 Objekte) bool {
	if !o.Metadatei.Equals(o1.Metadatei) {
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
	return o.Metadatei.AkteSha
}

func (o *Objekte) SetAkteSha(v schnittstellen.Sha) {
	o.Metadatei.AkteSha = sha.Make(v)
}
