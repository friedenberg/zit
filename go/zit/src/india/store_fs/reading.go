package store_fs

import (
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

// func (s *Store) ReadOneCheckedOut(
// 	o *sku.ObjekteOptions,
// 	em *KennungFDPair,
// ) (co *CheckedOut, err error) {
// 	co = GetCheckedOutPool().Get()

// 	if err = s.ReadOneInto(&em.Kennung, &co.Internal); err != nil {
// 		if collections.IsErrNotFound(err) {
// 			// TODO mark status as new
// 			err = nil
// 		} else {
// 			err = errors.Wrap(err)
// 			return
// 		}
// 	}

// 	if err = s.ReadOneExternalInto(
// 		o,
// 		em,
// 		&co.Internal,
// 		&co.External,
// 	); err != nil {
// 		if errors.Is(err, ErrExternalHasConflictMarker) {
// 			err = nil
// 			co.State = checked_out_state.StateConflicted
// 		} else {
// 			err = errors.Wrap(err)
// 		}

// 		return
// 	}

// 	// TODO move upstairs
// 	if err = s.tryRealizeAndOrStore(
// 		&co.External.Transacted,
// 		o,
// 	); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	co.DetermineState(false)

// 	return
// }

func (s *Store) ReadOneExternal(
	o *sku.ObjekteOptions,
	em *KennungFDPair,
	t *sku.Transacted,
) (e *External, err error) {
	e = GetExternalPool().Get()

	if err = s.ReadOneExternalInto(o, em, t, e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) UpdateTransacted(z *sku.Transacted) (err error) {
	e, ok := s.Get(&z.Kennung)

	if !ok {
		return
	}

	var e2 *External

	if e2, err = s.ReadExternalFromKennungFDPair(
		sku.ObjekteOptions{
			Mode: objekte_mode.ModeUpdateTai,
		},
		e,
		z,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = z.SetFromSkuLike(&e2.Transacted); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadOneExternalInto(
	o *sku.ObjekteOptions,
	em *KennungFDPair,
	t *sku.Transacted,
	e *External,
) (err error) {
	if err = e.ResetWithExternalMaybe(em); err != nil {
		err = errors.Wrap(err)
		return
	}

	var m checkout_mode.Mode

	if m, err = em.GetFDs().GetCheckoutModeOrError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = em.FDs.ConflictMarkerError(); err != nil {
		return
	}

	var t1 *sku.Transacted

	if t != nil {
		t1 = t
	}

	switch m {
	case checkout_mode.ModeAkteOnly:
		if err = s.ReadOneExternalAkte(e, t1); err != nil {
			err = errors.Wrap(err)
			return
		}

	case checkout_mode.ModeObjekteOnly, checkout_mode.ModeObjekteAndAkte:
		if e.FDs.Objekte.IsStdin() {
			if err = s.ReadOneExternalObjekteReader(os.Stdin, e); err != nil {
				err = errors.Wrap(err)
				return
			}
		} else {
			if err = s.ReadOneExternalObjekte(e, t1); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

	default:
		panic(checkout_mode.MakeErrInvalidCheckoutModeMode(m))
	}

	if !e.FDs.Akte.IsEmpty() {
		aFD := &e.FDs.Akte
		ext := aFD.ExtSansDot()
		typFromExtension := s.konfig.GetTypStringFromExtension(ext)

		if typFromExtension == "" {
			ui.Err().Printf("typ extension unknown: %s", aFD.ExtSansDot())
			typFromExtension = ext
		}

		if err = e.Transacted.Metadatei.Typ.Set(typFromExtension); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if o.Clock == nil {
		o.Clock = &e.FDs
	}

	return
}

func (s *Store) ReadOneExternalObjekte(
	e *External,
	t *sku.Transacted,
) (err error) {
	if t != nil {
		metadatei.Resetter.ResetWith(e.GetMetadatei(), t.GetMetadatei())
	}

	var f *os.File

	if f, err = files.Open(e.GetObjekteFD().GetPath()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	if err = s.ReadOneExternalObjekteReader(f, e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadOneExternalObjekteReader(
	r io.Reader,
	e *External,
) (err error) {
	if _, err = s.metadateiTextParser.ParseMetadatei(r, e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadOneExternalAkte(
	e *External,
	t *sku.Transacted,
) (err error) {
	metadatei.Resetter.ResetWith(&e.Metadatei, t.GetMetadatei())

	var aw sha.WriteCloser

	if aw, err = s.standort.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	var f *os.File

	if f, err = files.OpenExclusiveReadOnly(
		e.GetAkteFD().GetPath(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	if _, err = io.Copy(aw, f); err != nil {
		err = errors.Wrap(err)
		return
	}

	e.GetMetadatei().Akte.SetShaLike(aw)

	return
}
