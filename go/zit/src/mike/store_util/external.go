package store_util

import (
	"io"
	"os"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/src/bravo/files"
	"code.linenisgreat.com/zit/src/charlie/sha"
	"code.linenisgreat.com/zit/src/delta/thyme"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type ExternalReader interface {
	ReadOneExternal(
		em *sku.ExternalMaybe,
		t *sku.Transacted,
	) (e *sku.External, err error)

	ReadOneExternalObjekte(
		e *sku.External,
		t *sku.Transacted,
	) (err error)

	ReadOneExternalObjekteReader(
		r io.Reader,
		e *sku.External,
	) (err error)

	ReadOneExternalAkte(
		e *sku.External,
		t *sku.Transacted,
	) (err error)
}

func (s *common) ReadOneExternal(
	em *sku.ExternalMaybe,
	t *sku.Transacted,
) (e *sku.External, err error) {
	if err = em.FDs.ConflictMarkerError(); err != nil {
		return
	}

	var m checkout_mode.Mode

	if m, err = em.GetFDs().GetCheckoutModeOrError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	e = sku.GetExternalPool().Get()

	if err = e.ResetWithExternalMaybe(em); err != nil {
		err = errors.Wrap(err)
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
		if err = s.ReadOneExternalObjekte(e, t1); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		panic(checkout_mode.MakeErrInvalidCheckoutModeMode(m))
	}

	return
}

func (s *common) ReadOneExternalObjekte(
	e *sku.External,
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

func (s *common) ReadOneExternalObjekteReader(
	r io.Reader,
	e *sku.External,
) (err error) {
	if _, err = s.metadateiTextParser.ParseMetadatei(r, e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *common) ReadOneExternalAkte(
	e *sku.External,
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

	var fStat os.FileInfo

	if fStat, err = f.Stat(); err != nil {
		err = errors.Wrap(err)
		return
	}

	e.Metadatei.Tai = kennung.TaiFromTime(thyme.Tyme(fStat.ModTime()))
	e.GetMetadatei().Akte.SetShaLike(aw)

	if err = sku.CalculateAndSetSha(
		e,
		s.persistentMetadateiFormat,
		s.options,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
