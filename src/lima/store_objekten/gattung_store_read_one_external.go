package store_objekten

import (
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/checkout_mode"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/hotel/sku"
)

func (s *commonStore[O, OPtr, K, KPtr]) ReadOneExternal(
	em *sku.ExternalMaybe[K, KPtr],
	t *sku.Transacted[K, KPtr],
) (e sku.External[K, KPtr], err error) {
	var m checkout_mode.Mode

	if m, err = em.GetFDs().GetCheckoutMode(); err != nil {
		err = errors.Wrap(err)
		return
	}

	e.ResetWithExternalMaybe(*em)

	switch m {
	case checkout_mode.ModeAkteOnly:
		if err = s.ReadOneExternalAkte(&e, t); err != nil {
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

func (s *commonStore[O, OPtr, K, KPtr]) readOneExternalObjekte(
	e *sku.External[K, KPtr],
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

	return
}
