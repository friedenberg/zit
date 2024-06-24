package sku

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
)

type ExternalFS struct {
	Transacted
	FDs FDPair
}

func (t *ExternalFS) GetSkuExternalLike() ExternalLike {
	return t
}

func (c *ExternalFS) GetSku() *Transacted {
	return &c.Transacted
}

func (t *ExternalFS) SetFromSkuLike(sk SkuLike) (err error) {
	switch skt := sk.(type) {
	case *ExternalFS:
		t.FDs.ResetWith(skt.GetFDs())
	}

	if err = t.Transacted.SetFromSkuLike(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *ExternalFS) GetKennung() kennung.Kennung {
	return &a.Kennung
}

func (a *ExternalFS) GetMetadatei() *metadatei.Metadatei {
	return &a.Metadatei
}

func (a *ExternalFS) GetGattung() schnittstellen.GattungLike {
	return a.Kennung.GetGattung()
}

func (a *ExternalFS) String() string {
	return fmt.Sprintf(
		". %s %s %s %s",
		a.GetGattung(),
		a.GetKennung(),
		a.GetObjekteSha(),
		a.GetAkteSha(),
	)
}

func (a *ExternalFS) GetAkteSha() schnittstellen.ShaLike {
	return &a.Metadatei.Akte
}

func (a *ExternalFS) SetAkteSha(v schnittstellen.ShaLike) (err error) {
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

func (a *ExternalFS) AsTransacted() (b Transacted) {
	b = a.Transacted

	return
}

func (a *ExternalFS) GetFDs() *FDPair {
	return &a.FDs
}

func (a *ExternalFS) GetFDsPtr() *FDPair {
	return &a.FDs
}

func (a *ExternalFS) GetAkteFD() *fd.FD {
	return &a.FDs.Akte
}

func (a *ExternalFS) SetAkteFD(v *fd.FD) {
	a.FDs.Akte.ResetWith(v)
	a.Metadatei.Akte.SetShaLike(v.GetShaLike())
}

func (a *ExternalFS) GetAktePath() string {
	return a.FDs.Akte.GetPath()
}

func (a *ExternalFS) GetObjekteFD() *fd.FD {
	return &a.FDs.Objekte
}

func (a *ExternalFS) ResetWithExternalMaybe(
	b *KennungFDPair,
) (err error) {
	k := b.GetKennungLike()
	a.Kennung.ResetWithKennung(k)
	metadatei.Resetter.Reset(&a.Metadatei)
	a.FDs.ResetWith(b.GetFDs())

	return
}

func (o *ExternalFS) GetKey() string {
	return fmt.Sprintf("%s.%s", o.GetGattung(), o.GetKennung())
}

func (e *ExternalFS) GetCheckoutMode() (m checkout_mode.Mode, err error) {
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

func (lessorExternal) Less(a, b ExternalFS) bool {
	panic("not supported")
}

func (lessorExternal) LessPtr(a, b *ExternalFS) bool {
	return a.GetTai().Less(b.GetTai())
}

type equalerExternal struct{}

func (equalerExternal) Equals(a, b ExternalFS) bool {
	panic("not supported")
}

func (equalerExternal) EqualsPtr(a, b *ExternalFS) bool {
	return a.EqualsSkuLikePtr(b)
}
