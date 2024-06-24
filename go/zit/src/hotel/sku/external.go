package sku

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
)

type External struct {
	Transacted
}

func (t *External) GetSkuExternalLike() ExternalLike {
	return t
}

func (c *External) GetSku() *Transacted {
	return &c.Transacted
}

func (t *External) SetFromSkuLike(sk SkuLike) (err error) {
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

	return
}

func (a *External) AsTransacted() (b Transacted) {
	b = a.Transacted

	return
}

func (a *External) ResetWithExternalMaybe(
	b *KennungFDPair,
) (err error) {
	k := b.GetKennungLike()
	a.Kennung.ResetWithKennung(k)
	metadatei.Resetter.Reset(&a.Metadatei)

	return
}

func (o *External) GetKey() string {
	return fmt.Sprintf("%s.%s", o.GetGattung(), o.GetKennung())
}

// func (e *External) GetCheckoutMode() (m checkout_mode.Mode, err error) {
// 	switch {
// 	case !e.FDs.Objekte.IsEmpty() && !e.FDs.Akte.IsEmpty():
// 		m = checkout_mode.ModeObjekteAndAkte

// 	case !e.FDs.Akte.IsEmpty():
// 		m = checkout_mode.ModeAkteOnly

// 	case !e.FDs.Objekte.IsEmpty():
// 		m = checkout_mode.ModeObjekteOnly

// 	default:
// 		err = checkout_mode.MakeErrInvalidCheckoutMode(
// 			errors.Errorf("all FD's are empty"),
// 		)
// 	}

// 	return
// }
