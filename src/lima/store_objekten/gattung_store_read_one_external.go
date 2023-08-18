package store_objekten

import (
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/checkout_mode"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/hotel/objekte"
)

func (s *commonStore[O, OPtr, K, KPtr]) ReadOneExternal(
	em *sku.ExternalMaybe[K, KPtr],
	t *sku.Transacted[K, KPtr],
) (e objekte.External[O, OPtr, K, KPtr], err error) {
	var m checkout_mode.Mode

	if m, err = em.GetFDs().GetCheckoutMode(); err != nil {
		err = errors.Wrap(err)
		return
	}

	e.Sku.ResetWithExternalMaybe(*em)

	switch m {
	case checkout_mode.ModeAkteOnly:
		if err = s.readOneExternalAkte(&e, t); err != nil {
			err = errors.Wrap(err)
			return
		}

	case checkout_mode.ModeObjekteOnly, checkout_mode.ModeObjekteAndAkte:
		if err = s.readOneExternalObjekte(&e, t); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *commonStore[O, OPtr, K, KPtr]) readOneExternalAkte(
	e *objekte.External[O, OPtr, K, KPtr],
	t *sku.Transacted[K, KPtr],
) (err error) {
	e.SetMetadatei(t.GetMetadatei())

	var aw sha.WriteCloser

	if aw, err = s.StoreUtil.AkteWriter(); err != nil {
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
	e.SetAkteSha(sh)

	if err = s.SaveObjekte(e); err != nil {
		err = errors.Wrapf(err, "%s", f.Name())
		return
	}

	return
}

func (s *commonStore[O, OPtr, K, KPtr]) readOneExternalObjekte(
	e *objekte.External[O, OPtr, K, KPtr],
	t *sku.Transacted[K, KPtr],
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

	if _, err = s.textParser.ParseMetadatei(f, e); err != nil {
		err = errors.Wrapf(err, "%s", f.Name())
		return
	}

	if err = s.SaveObjekte(e); err != nil {
		err = errors.Wrapf(err, "%s", f.Name())
		return
	}

	return
}
