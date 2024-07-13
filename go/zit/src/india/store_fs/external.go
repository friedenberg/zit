package store_fs

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type External struct {
	sku.Transacted
	FDs FDPair
}

func (t *External) GetSkuExternalLike() sku.ExternalLike {
	return t
}

func (c *External) GetSku() *sku.Transacted {
	return &c.Transacted
}

func (t *External) SetFromSkuLike(sk sku.SkuLike) (err error) {
	switch skt := sk.(type) {
	case *External:
		t.FDs.ResetWith(skt.GetFDs())
	}

	if err = t.Transacted.SetFromSkuLike(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *External) GetKennung() kennung.IdLike {
	return &a.Kennung
}

func (a *External) GetMetadatei() *metadatei.Metadatei {
	return &a.Metadatei
}

func (a *External) GetGattung() interfaces.Genre {
	return a.Kennung.GetGenre()
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

func (a *External) GetAkteSha() interfaces.ShaLike {
	return &a.Metadatei.Akte
}

func (a *External) SetAkteSha(v interfaces.ShaLike) (err error) {
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

func (a *External) GetFDs() *FDPair {
	return &a.FDs
}

func (a *External) GetFDsPtr() *FDPair {
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
	b *KennungFDPair,
) (err error) {
	k := b.GetKennungLike()
	a.Kennung.ResetWithIdLike(k)
	metadatei.Resetter.Reset(&a.Metadatei)
	a.FDs.ResetWith(b.GetFDs())

	return
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
