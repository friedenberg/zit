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

type External struct {
	Transacted
	FDs ExternalFDs
}

func (t *External) SetFromSkuLike(sk SkuLikePtr) (err error) {
	switch skt := sk.(type) {
	case SkuLikeExternalPtr:
		t.FDs = skt.GetFDs()
	}

	if err = t.Transacted.SetFromSkuLike(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a External) GetKennung() kennung.Kennung {
	return a.Kennung
}

func (a External) GetMetadatei() metadatei.Metadatei {
	return a.Metadatei
}

func (a *External) GetMetadateiPtr() *metadatei.Metadatei {
	return &a.Metadatei
}

func (a External) GetGattung() schnittstellen.GattungLike {
	return a.Kennung.GetGattung()
}

func (a External) GetKennungLike() kennung.Kennung {
	return a.Kennung
}

func (a External) String() string {
	return fmt.Sprintf(
		". %s %s %s %s",
		a.GetGattung(),
		a.GetKennung(),
		a.ObjekteSha,
		a.GetAkteSha(),
	)
}

func (a External) GetAkteSha() schnittstellen.ShaLike {
	return a.Metadatei.AkteSha
}

func (a *External) SetAkteSha(v schnittstellen.ShaLike) {
	sh := sha.Make(v)
	a.Metadatei.AkteSha = sh
	a.FDs.Akte.Sha = sh
}

func (a *External) AsTransacted() (b Transacted) {
	b = a.Transacted

	return
}

func (a External) GetFDs() ExternalFDs {
	return a.FDs
}

func (a *External) GetFDsPtr() *ExternalFDs {
	return &a.FDs
}

func (a External) GetAkteFD() kennung.FD {
	return a.FDs.Akte
}

func (a External) GetAktePath() string {
	return a.FDs.Akte.GetPath()
}

func (a External) GetObjekteFD() kennung.FD {
	return a.FDs.Objekte
}

func (a *External) Reset() {
	a.ObjekteSha.Reset()
	a.Kennung.Reset()
	a.Metadatei.Reset()
}

func (a *External) ResetWith(b *External) {
	a.ObjekteSha.ResetWith(b.ObjekteSha)
	a.Kennung.ResetWithKennung(b.Kennung)
	a.Metadatei.ResetWith(b.GetMetadatei())
}

func (a *External) ResetWithExternalMaybe(
	b ExternalMaybe,
) (err error) {
	k := b.GetKennungLike()
	a.Kennung.ResetWithKennung(k)
	a.ObjekteSha.Reset()
	a.Metadatei.Reset()
	a.FDs = b.GetFDs()

	return
}

func (a External) EqualsAny(b any) (ok bool) {
	return values.Equals(a, b)
}

func (a External) EqualsSkuLikePtr(b SkuLikePtr) (ok bool) {
	return values.Equals(a, b)
}

func (a External) Equals(b External) (ok bool) {
	if !kennung.Equals(a.GetKennung(), b.GetKennung()) {
		return
	}

	if !a.ObjekteSha.Equals(b.ObjekteSha) {
		return
	}

	return true
}

func (o External) GetKey() string {
	return fmt.Sprintf("%s.%s", o.GetGattung(), o.GetKennung())
}

func (e External) GetCheckoutMode() (m checkout_mode.Mode, err error) {
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
