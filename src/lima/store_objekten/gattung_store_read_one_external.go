package store_objekten

import (
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/checkout_mode"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
)

func (s *commonStore[O, OPtr, K, KPtr, V, VPtr]) ReadOneExternal(
	em sku.ExternalMaybe[K, KPtr],
	t *objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
) (e objekte.External[O, OPtr, K, KPtr], err error) {
	var m checkout_mode.Mode

	if m, err = em.GetCheckoutMode(); err != nil {
		err = errors.Wrap(err)
		return
	}

	e.Sku.ResetWithExternalMaybe(em)

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

	objekte.CorrectAkteShaWith(&e, e)

	var ar sha.ReadCloser

	if ar, err = s.AkteReader(e.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ar)

	var akteSha schnittstellen.Sha

	if akteSha, _, err = s.AkteFormat.ParseSaveAkte(ar, &e.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !akteSha.EqualsSha(e.GetAkteSha()) {
		panic(errors.Errorf(
			"akte sha in parser was %s while akte sha in objekte was %s",
			akteSha,
			e.GetAkteSha(),
		))
	}

	return
}

func (s *commonStore[O, OPtr, K, KPtr, V, VPtr]) readOneExternalAkte(
	e *objekte.External[O, OPtr, K, KPtr],
	t *objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
) (err error) {
	e.Objekte = t.Akte
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

	sh := sha.Make(aw.Sha())
	e.SetAkteSha(sh)

	if err = s.SaveObjekte(e); err != nil {
		err = errors.Wrapf(err, "%s", f.Name())
		return
	}

	return
}

func (s *commonStore[O, OPtr, K, KPtr, V, VPtr]) readOneExternalObjekte(
	e *objekte.External[O, OPtr, K, KPtr],
	t *objekte.Transacted[O, OPtr, K, KPtr, V, VPtr],
) (err error) {
	var f *os.File

	if f, err = files.Open(e.GetObjekteFD().Path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	if t != nil {
		OPtr(&e.Objekte).ResetWith(t.Akte)
		e.Metadatei.ResetWith(t.Metadatei)
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
