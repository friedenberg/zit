package sku

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/checkout_mode"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type External[T kennung.KennungLike[T], T1 kennung.KennungLikePtr[T]] struct {
	ObjekteSha sha.Sha
	AkteSha    sha.Sha
	// TODO-P4 turn in ExternalMaybe
	Kennung T
	FDs     ExternalFDs
}

func (a External[T, T1]) String() string {
	return fmt.Sprintf(
		". %s %s %s %s",
		a.Kennung.GetGattung(),
		a.Kennung,
		a.ObjekteSha,
		a.AkteSha,
	)
}

func (a External[T, T1]) GetAkteSha() schnittstellen.Sha {
	return a.AkteSha
}

func (a *External[T, T1]) Transacted() (b Transacted[T, T1]) {
	b = Transacted[T, T1]{
		Kennung:    a.Kennung,
		ObjekteSha: a.ObjekteSha,
		AkteSha:    a.AkteSha,
	}

	return
}

func (a *External[T, T1]) Reset() {
	a.ObjekteSha = sha.Sha{}
	a.AkteSha = sha.Sha{}
	T1(&a.Kennung).Reset()
}

func (a *External[T, T1]) ResetWith(b *External[T, T1]) {
	a.ObjekteSha = b.ObjekteSha
	a.AkteSha = b.AkteSha
	T1(&a.Kennung).ResetWith(b.Kennung)
}

func (a *External[T, T1]) ResetWithExternalMaybe(b ExternalMaybe[T, T1]) {
	todo.Change("use this in other places")
	a.ObjekteSha = sha.Sha{}
	a.AkteSha = sha.Sha{}
	a.FDs = b.FDs
	T1(&a.Kennung).ResetWith(b.Kennung)
}

func (a External[T, T1]) Equals(b External[T, T1]) (ok bool) {
	if a.Kennung.Equals(b.Kennung) {
		return
	}

	if !a.ObjekteSha.Equals(b.ObjekteSha) {
		return
	}

	return true
}

func (o External[T, T1]) GetKey() string {
	return fmt.Sprintf("%s.%s", o.Kennung.GetGattung(), o.Kennung)
}

func (e External[T, T1]) GetCheckoutMode() (m checkout_mode.Mode, err error) {
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
