package sku

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/checkout_mode"
	"github.com/friedenberg/zit/src/echo/fd"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type External struct {
	Transacted
	FDs ExternalFDs
}

func (t *External) GetExternalSkuLike() SkuExternalLike {
	return t
}

func (t *External) SetFromSkuLike(sk SkuLike) (err error) {
	switch skt := sk.(type) {
	case SkuExternalLike:
		t.FDs.ResetWith(skt.GetFDs())
	}

	if err = t.Transacted.SetFromSkuLike(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *External) GetKennung() kennung.Kennung {
	return &a.Kennung
}

func (a *External) GetMetadatei() *metadatei.Metadatei {
	return &a.Metadatei
}

func (a *External) GetGattung() schnittstellen.GattungLike {
	return a.Kennung.GetGattung()
}

func (a *External) String() string {
	return fmt.Sprintf(
		". %s %s %s %s",
		a.GetGattung(),
		a.GetKennung(),
		a.GetObjekteSha(),
		a.GetAkteSha(),
	)
}

func (a *External) GetAkteSha() schnittstellen.ShaLike {
	return &a.Metadatei.Akte
}

func (a *External) SetAkteSha(v schnittstellen.ShaLike) (err error) {
	if err = a.Metadatei.Akte.SetShaLike(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = a.FDs.Akte.SetShaLike(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *External) AsTransacted() (b Transacted) {
	b = a.Transacted

	return
}

func (a *External) GetFDs() *ExternalFDs {
	return &a.FDs
}

func (a *External) GetFDsPtr() *ExternalFDs {
	return &a.FDs
}

func (a *External) GetAkteFD() *fd.FD {
	return &a.FDs.Akte
}

func (a *External) SetAkteFD(v *fd.FD) {
	a.FDs.Akte.ResetWith(v)
	a.Metadatei.Akte.SetShaLike(v.GetShaLike())
}

func (a *External) GetAktePath() string {
	return a.FDs.Akte.GetPath()
}

func (a *External) GetObjekteFD() *fd.FD {
	return &a.FDs.Objekte
}

func (a *External) ResetWithExternalMaybe(
	b *ExternalMaybe,
) (err error) {
	k := b.GetKennungLike()
	a.Kennung.ResetWithKennung(k)
	metadatei.Resetter.Reset(&a.Metadatei)
	a.FDs.ResetWith(b.GetFDs())

	return
}

func (a *External) EqualsSkuLikePtr(b SkuLike) (ok bool) {
	if !kennung.Equals(a.GetKennung(), b.GetKennung()) {
		return
	}

	if !a.GetObjekteSha().EqualsSha(b.GetObjekteSha()) {
		return
	}

	return true
}

func (o *External) GetKey() string {
	return fmt.Sprintf("%s.%s", o.GetGattung(), o.GetKennung())
}

func (e *External) GetCheckoutMode() (m checkout_mode.Mode, err error) {
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

type lessorExternal struct{}

func (lessorExternal) Less(a, b External) bool {
	panic("not supported")
}

func (lessorExternal) LessPtr(a, b *External) bool {
	return a.GetTai().Less(b.GetTai())
}

type equalerExternal struct{}

func (equalerExternal) Equals(a, b External) bool {
	panic("not supported")
}

func (equalerExternal) EqualsPtr(a, b *External) bool {
	return a.EqualsSkuLikePtr(b)
}
