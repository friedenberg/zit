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

type External[K kennung.KennungLike[K], KPtr kennung.KennungLikePtr[K]] struct {
	Transacted[K, KPtr]
	FDs ExternalFDs
}

func (a External[K, KPtr]) GetKennung() K {
	return a.Kennung
}

func (a External[K, KPtr]) GetMetadatei() metadatei.Metadatei {
	return a.Metadatei
}

func (a *External[K, KPtr]) GetMetadateiPtr() *metadatei.Metadatei {
	return &a.Metadatei
}

func (a External[K, KPtr]) GetGattung() schnittstellen.GattungLike {
	return a.Kennung.GetGattung()
}

func (a External[K, KPtr]) GetKennungLike() kennung.Kennung {
	return a.Kennung
}

func (a External[K, KPtr]) String() string {
	return fmt.Sprintf(
		". %s %s %s %s",
		a.GetGattung(),
		a.GetKennung(),
		a.ObjekteSha,
		a.GetAkteSha(),
	)
}

func (a External[K, KPtr]) GetAkteSha() schnittstellen.ShaLike {
	return a.Metadatei.AkteSha
}

func (a *External[K, KPtr]) SetAkteSha(v schnittstellen.ShaLike) {
	sh := sha.Make(v)
	a.Metadatei.AkteSha = sh
	a.FDs.Akte.Sha = sh
}

func (a *External[K, KPtr]) AsTransacted() (b Transacted[K, KPtr]) {
	b = Transacted[K, KPtr]{
		Kennung: a.GetKennung(),
		Metadatei: metadatei.Metadatei{
			AkteSha: sha.Make(a.GetAkteSha()),
		},
		ObjekteSha: a.ObjekteSha,
	}

	return
}

func (a External[K, KPtr]) GetFDs() ExternalFDs {
	return a.FDs
}

func (a *External[K, KPtr]) GetFDsPtr() *ExternalFDs {
	return &a.FDs
}

func (a External[K, KPtr]) GetAkteFD() kennung.FD {
	return a.FDs.Akte
}

func (a External[K, KPtr]) GetAktePath() string {
	return a.FDs.Akte.Path
}

func (a External[K, KPtr]) GetObjekteFD() kennung.FD {
	return a.FDs.Objekte
}

func (a *External[K, KPtr]) Reset() {
	a.ObjekteSha.Reset()
	KPtr(&a.Kennung).Reset()
	a.Metadatei.Reset()
}

func (a *External[K, KPtr]) ResetWith(b *External[K, KPtr]) {
	a.ObjekteSha.ResetWith(b.ObjekteSha)
	KPtr(&a.Kennung).ResetWith(b.Kennung)
	a.Metadatei.ResetWith(b.GetMetadatei())
}

func (a *External[K, KPtr]) ResetWithExternalMaybe(
	b ExternalMaybe,
) (err error) {
	k := b.GetKennungLike()

	switch kt := k.(type) {
	case K:
		KPtr(&a.Kennung).ResetWith(kt)

	case KPtr:
		KPtr(&a.Kennung).ResetWith(K(*kt))

	default:
		err = errors.Errorf("unsupported kennung type: %T", kt)
		return
	}

	a.ObjekteSha.Reset()
	a.Metadatei.Reset()
	a.FDs = b.GetFDs()

	return
}

func (a External[K, KPtr]) EqualsAny(b any) (ok bool) {
	return values.Equals(a, b)
}

func (a External[K, KPtr]) EqualsSkuLike(b SkuLike) (ok bool) {
	return values.Equals(a, b)
}

func (a External[K, KPtr]) Equals(b External[K, KPtr]) (ok bool) {
	if a.GetKennung().Equals(b.GetKennung()) {
		return
	}

	if !a.ObjekteSha.Equals(b.ObjekteSha) {
		return
	}

	return true
}

func (o External[K, KPtr]) GetKey() string {
	return fmt.Sprintf("%s.%s", o.GetGattung(), o.GetKennung())
}

func (e External[K, KPtr]) GetCheckoutMode() (m checkout_mode.Mode, err error) {
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
