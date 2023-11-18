package sku

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/golf/objekte_format"
)

func CalculateAndConfirmSha(
	sk SkuLike,
	format objekte_format.Formatter,
	o objekte_format.Options,
	sh schnittstellen.ShaLike,
) (err error) {
	if err = CalculateAndSetSha(sk, format, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !sk.GetObjekteSha().EqualsSha(sh) {
		err = errors.Errorf(
			"expected sha %q but got %q",
			sh,
			sk.GetObjekteSha(),
		)

		err = errors.Wrapf(err, "Format: %T", format)

		return
	}

	return
}

func CalculateAndSetSha(
	sk SkuLike,
	format objekte_format.Formatter,
	o objekte_format.Options,
) (err error) {
	w := sha.MakeWriter(nil)

	if _, err = format.FormatPersistentMetadatei(w, sk, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh := w.GetShaLike()
	sk.SetObjekteSha(sh)

	return
}

func ReadFromSha(
	sk SkuLike,
	orf schnittstellen.ObjekteReaderFactory,
	format objekte_format.Format,
	o objekte_format.Options,
) (err error) {
	expected := sk.GetObjekteSha()

	var or schnittstellen.ShaReadCloser

	if or, err = orf.ObjekteReader(expected); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, or)

	if _, err = format.ParsePersistentMetadatei(or, sk, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	actual := or.GetShaLike()

	if !actual.EqualsSha(expected) {
		err = errors.Errorf("expected %s but got %s", expected, actual)
		return
	}

	return
}
