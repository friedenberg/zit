package store_util

import (
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/checkout_mode"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/hotel/sku"
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

	ReadOneExternalAkte(
		e *sku.External,
		t *sku.Transacted,
	) (err error)
}

func (s *common) ReadOneExternal(
	em *sku.ExternalMaybe,
	t *sku.Transacted,
) (e *sku.External, err error) {
	var m checkout_mode.Mode

	if m, err = em.GetFDs().GetCheckoutMode(); err != nil {
		err = errors.Wrap(err)
		return
	}

	e = &sku.External{}

	if err = e.ResetWithExternalMaybe(*em); err != nil {
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
	}

	return
}

func (s *common) ReadOneExternalObjekte(
	e *sku.External,
	t *sku.Transacted,
) (err error) {
	var f *os.File

	if f, err = files.Open(e.GetObjekteFD().Path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	if t != nil {
		e.GetMetadateiPtr().ResetWith(t.GetMetadatei())
	}

	if _, err = s.metadateiTextParser.ParseMetadatei(f, e); err != nil {
		err = errors.Wrapf(err, "%s", f.Name())
		return
	}

	return
}

func (s *common) ReadOneExternalAkte(
	e *sku.External,
	t *sku.Transacted,
) (err error) {
	e.SetMetadatei(t.GetMetadatei())

	var aw sha.WriteCloser

	if aw, err = s.standort.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	var f *os.File

	if f, err = files.OpenExclusiveReadOnly(
		e.GetAkteFD().Path,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	if _, err = io.Copy(aw, f); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh := sha.Make(aw.GetShaLike())
	e.GetMetadateiPtr().AkteSha = sh

	return
}
