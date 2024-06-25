package store_fs

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
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
	a.Kennung.ResetWithKennung(k)
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

func (s *Store) CombineOneCheckedOutFS(
	sk2 *sku.Transacted,
) (co *CheckedOut, err error) {
	co = GetCheckedOutPool().Get()

	if err = co.Internal.SetFromSkuLike(sk2); err != nil {
		err = errors.Wrap(err)
		return
	}

	ok := false

	var e *KennungFDPair

	if e, ok = s.Get(&sk2.Kennung); !ok {
		err = iter.MakeErrStopIteration()
		return
	}

	var e2 *External

	if e2, err = s.ReadOneExternalFS(
		sku.ObjekteOptions{
			Mode: objekte_mode.ModeUpdateTai,
		},
		e,
		sk2,
	); err != nil {
		if errors.IsNotExist(err) {
			err = iter.MakeErrStopIteration()
		} else if errors.Is(err, ErrExternalHasConflictMarker) {
			co.State = checked_out_state.StateConflicted
			co.External.FDs = e.FDs

			if err = co.External.Kennung.SetWithKennung(&sk2.Kennung); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		} else {
			err = errors.Wrapf(err, "Cwd: %#v", e)
		}

		return
	}

	if err = co.External.SetFromSkuLike(e2); err != nil {
		err = errors.Wrap(err)
		return
	}

	co.DetermineState(false)

	return
}

func (s *Store) ReadOneExternalFS(
	o sku.ObjekteOptions,
	em *KennungFDPair,
	t *sku.Transacted,
) (e *External, err error) {
	e = GetExternalPool().Get()

	if err = s.ReadOneExternalFSInto(o, em, t, e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadOneExternalFSInto(
	o sku.ObjekteOptions,
	em *KennungFDPair,
	t *sku.Transacted,
	e *External,
) (err error) {
	o.Del(objekte_mode.ModeApplyProto)

	if err = s.ReadOneExternalInto(
		&o,
		em,
		t,
		e,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.storeFuncs.FuncCommit(
		&e.Transacted,
		o,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
