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
	Metadatei metadatei.Metadatei
}

func (z *Objekte) SetAkteSha(v schnittstellen.Sha) {
	z.Metadatei.AkteSha = sha.Make(v)
}

func (z Objekte) GetMetadatei() metadatei.Metadatei {
	return z.Metadatei
}

func (z *Objekte) GetMetadateiPtr() *metadatei.Metadatei {
	return &z.Metadatei
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
	return z.Metadatei.AkteSha
}

func (z Objekte) Equals(z1 Objekte) bool {
	errors.TodoP4("figure out why this doesn't always work for `status`")
	if !z.Metadatei.AkteSha.Equals(z1.Metadatei.AkteSha) {
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

	if !z.Metadatei.AkteSha.IsNull() {
		return false
	}

	return true
}

func (z *Objekte) Reset() {
	z.Metadatei.Reset()
}

func (z *Objekte) ResetWith(z1 Objekte) {
	z.Metadatei.ResetWith(z1.Metadatei)
}
