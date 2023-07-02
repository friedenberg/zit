package sku

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/checkout_mode"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type External[K kennung.KennungLike[K], KPtr kennung.KennungLikePtr[K]] struct {
	ObjekteSha  sha.Sha
	WithKennung metadatei.WithKennung[K, KPtr]
	FDs         ExternalFDs
}

func (a External[K, KPtr]) GetKennung() K {
	return a.WithKennung.GetKennung()
}

func (a External[K, KPtr]) GetMetadatei() Metadatei {
	return a.WithKennung.GetMetadatei()
}

func (a *External[K, KPtr]) GetMetadateiPtr() *Metadatei {
	return a.WithKennung.GetMetadateiPtr()
}

func (a External[K, KPtr]) GetGattung() gattung.Gattung {
	return gattung.Must(a.WithKennung.GetGattung())
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
	return a.WithKennung.Metadatei.AkteSha
}

func (a *External[K, KPtr]) SetAkteSha(v schnittstellen.ShaLike) {
	sh := sha.Make(v)
	a.WithKennung.Metadatei.AkteSha = sh
	a.FDs.Akte.Sha = sh
}

func (a *External[K, KPtr]) Transacted() (b Transacted[K, KPtr]) {
	b = Transacted[K, KPtr]{
		WithKennung: metadatei.WithKennung[K, KPtr]{
			Kennung: a.GetKennung(),
			Metadatei: Metadatei{
				AkteSha: sha.Make(a.GetAkteSha()),
			},
		},
		ObjekteSha: a.ObjekteSha,
	}

	return
}

func (a *External[K, KPtr]) Reset() {
	a.ObjekteSha.Reset()
	a.WithKennung.Reset()
}

func (a *External[K, KPtr]) ResetWith(b *External[K, KPtr]) {
	a.ObjekteSha.ResetWith(b.ObjekteSha)
	a.WithKennung.ResetWith(b.WithKennung)
}

func (a *External[K, KPtr]) ResetWithExternalMaybe(b ExternalMaybe[K, KPtr]) {
	todo.Change("use this in other places")
	a.ObjekteSha.Reset()
	a.WithKennung.Metadatei.AkteSha.Reset()
	a.FDs = b.FDs
	KPtr(&a.WithKennung.Kennung).ResetWith(b.Kennung)
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
		err = MakeErrInvalidCheckoutMode(
			errors.Errorf("all FD's are empty"),
		)
	}

	return
}
