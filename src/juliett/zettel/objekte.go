package zettel

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type Objekte struct {
	Akte      sha.Sha
	Metadatei metadatei.Metadatei
}

func (z Objekte) GetTyp() kennung.Typ {
	return z.Metadatei.Typ
}

func (z Objekte) GetEtiketten() schnittstellen.Set[kennung.Etikett] {
	return z.Metadatei.GetEtiketten()
}

func (z Objekte) GetGattung() schnittstellen.Gattung {
	return gattung.Zettel
}

func (z Objekte) GetAkteSha() schnittstellen.Sha {
	return z.Akte
}

func (z *Objekte) SetAkteSha(v schnittstellen.Sha) {
	z.Akte = sha.Make(v)
}

func (z Objekte) Equals(z1 Objekte) bool {
	errors.TodoP4("figure out why this doesn't always work for `status`")
	if !z.Akte.Equals(z1.Akte) {
		return false
	}

	if !z.Metadatei.Equals(z1.Metadatei) {
		return false
	}

	return true
}

func (z Objekte) IsEmpty() bool {
	if !z.Metadatei.IsEmpty() {
		return false
	}

	if !z.Akte.IsNull() {
		return false
	}

	return true
}

func (z *Objekte) Reset() {
	z.Akte = sha.Sha{}
	z.Metadatei.Reset()
}

func (z *Objekte) ResetWith(z1 Objekte) {
	z.Akte = z1.Akte
	z.Metadatei.ResetWith(z1.Metadatei)
}
