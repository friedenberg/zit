package sku

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
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

func (a *External) GetKennung() ids.IdLike {
	return &a.Kennung
}

func (a *External) GetMetadatei() *object_metadata.Metadata {
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
		a.GetObjectSha(),
		a.GetAkteSha(),
	)
}

func (a *External) GetAkteSha() interfaces.Sha {
	return &a.Metadatei.Blob
}

func (a *External) SetAkteSha(v interfaces.Sha) (err error) {
	if err = a.Metadatei.Blob.SetShaLike(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *External) AsTransacted() (b Transacted) {
	b = a.Transacted

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
