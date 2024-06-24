package store

import (
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/store_fs"
)

type ObjekteOptions struct {
	objekte_mode.Mode
	kennung.Clock
}

func (s *Store) ReadOneExternal(
	o ObjekteOptions,
	k1 schnittstellen.StringerGattungKastenGetter,
	sk *sku.Transacted,
) (el sku.ExternalLike, err error) {
	switch k1.GetKasten().GetKastenString() {
	case "chrome":
		// TODO populate with chrome kasten
		ui.Debug().Print("would populate from chrome")

	default:
		e, ok := s.cwdFiles.Get(k1)

		if !ok {
			return
		}

		if el, err = s.ReadOneExternalFS(o, e, sk); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) ReadOneCheckedOutFS(
	o ObjekteOptions,
	em *sku.KennungFDPair,
) (co *store_fs.CheckedOut, err error) {
	co = store_fs.GetCheckedOutPool().Get()

	if err = s.ReadOneInto(&em.Kennung, &co.Internal); err != nil {
		if collections.IsErrNotFound(err) {
			// TODO mark status as new
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	if err = s.ReadOneExternalInto(
		o,
		em,
		&co.Internal,
		&co.External,
	); err != nil {
		if errors.Is(err, sku.ErrExternalHasConflictMarker) {
			err = nil
			co.State = checked_out_state.StateConflicted
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	co.DetermineState(false)

	return
}

func (s *Store) ReadOneExternalFS(
	o ObjekteOptions,
	em *sku.KennungFDPair,
	t *sku.Transacted,
) (e *store_fs.External, err error) {
	e = store_fs.GetExternalPool().Get()

	if err = s.ReadOneExternalInto(o, em, t, e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadOneExternalInto(
	o ObjekteOptions,
	em *sku.KennungFDPair,
	t *sku.Transacted,
	e *store_fs.External,
) (err error) {
	var m checkout_mode.Mode

	if m, err = em.GetFDs().GetCheckoutModeOrError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = e.ResetWithExternalMaybe(em); err != nil {
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

	if o.Clock == nil {
		o.Clock = &e.FDs
	}

	if err = s.tryRealizeAndOrStore(
		&e.Transacted,
		o,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadOneExternalObjekte(
	e *store_fs.External,
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
	e *store_fs.External,
) (err error) {
	if _, err = s.metadateiTextParser.ParseMetadatei(r, e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadOneExternalAkte(
	e *store_fs.External,
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
