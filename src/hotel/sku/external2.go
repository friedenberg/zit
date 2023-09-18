package sku

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/checkout_mode"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type External2 struct {
	Transacted2
	FDs ExternalFDs
}

func (a External2) GetKennung() kennung.Kennung {
	return a.Kennung
}

func (a External2) GetMetadatei() metadatei.Metadatei {
	return a.Metadatei
}

func (a *External2) GetMetadateiPtr() *metadatei.Metadatei {
	return &a.Metadatei
}

func (a External2) GetGattung() schnittstellen.GattungLike {
	return a.Kennung.GetGattung()
}

func (a External2) GetKennungLike() kennung.Kennung {
	return a.Kennung
}

func (a External2) String() string {
	return fmt.Sprintf(
		". %s %s %s %s",
		a.GetGattung(),
		a.GetKennung(),
		a.ObjekteSha,
		a.GetAkteSha(),
	)
}

func (a External2) GetAkteSha() schnittstellen.ShaLike {
	return a.Metadatei.AkteSha
}

func (a *External2) SetAkteSha(v schnittstellen.ShaLike) {
	sh := sha.Make(v)
	a.Metadatei.AkteSha = sh
	a.FDs.Akte.Sha = sh
}

func (a *External2) AsTransacted() (b Transacted2) {
	b = a.Transacted2

	return
}

func (a External2) GetFDs() ExternalFDs {
	return a.FDs
}

func (a *External2) GetFDsPtr() *ExternalFDs {
	return &a.FDs
}

func (a External2) GetAkteFD() kennung.FD {
	return a.FDs.Akte
}

func (a External2) GetAktePath() string {
	return a.FDs.Akte.Path
}

func (a External2) GetObjekteFD() kennung.FD {
	return a.FDs.Objekte
}

func (a *External2) Reset() {
	a.ObjekteSha.Reset()
	a.Kennung.Reset()
	a.Metadatei.Reset()
}

func (a *External2) ResetWith(b *External2) {
	a.ObjekteSha.ResetWith(b.ObjekteSha)
	a.Kennung.ResetWithKennung(b.Kennung)
	a.Metadatei.ResetWith(b.GetMetadatei())
}

func (a *External2) ResetWithExternalMaybe(
	b ExternalMaybe,
) (err error) {
	k := b.GetKennungLike()
	a.Kennung.ResetWithKennung(k)
	a.ObjekteSha.Reset()
	a.Metadatei.Reset()
	a.FDs = b.GetFDs()

	return
}

func (a External2) EqualsAny(b any) (ok bool) {
	return values.Equals(a, b)
}

func (a External2) EqualsSkuLike(b SkuLike) (ok bool) {
	return values.Equals(a, b)
}

func (a External2) Equals(b External2) (ok bool) {
	if !kennung.Equals(a.GetKennung(), b.GetKennung()) {
		return
	}

	if !a.ObjekteSha.Equals(b.ObjekteSha) {
		return
	}

	return true
}

func (o External2) GetKey() string {
	return fmt.Sprintf("%s.%s", o.GetGattung(), o.GetKennung())
}

func (e External2) GetCheckoutMode() (m checkout_mode.Mode, err error) {
	switch {
	case !e.FDs.Objekte.IsEmpty() && !e.FDs.Akte.IsEmpty():
		m = checkout_mode.ModeObjekteAndAkte

	case !e.FDs.Akte.IsEmpty():
		m = checkout_mode.ModeAkteOnly

	case !e.FDs.Objekte.IsEmpty():
		m = checkout_mode.ModeObjekteOnly

	default:
		err = checkout_mode.MakeErrInvalidCheckoutMode(
			errors.Errorf("all FD's are empty"),
		)
	}

	return
}
