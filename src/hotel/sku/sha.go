package sku

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/golf/objekte_format"
)

func CalculateAndConfirmSha(
	sk SkuLikePtr,
	format objekte_format.Format,
	sh schnittstellen.ShaLike,
) (err error) {
	if err = CalculateAndSetSha(sk, format); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !sk.GetObjekteSha().EqualsSha(sh) {
		err = errors.Errorf(
			"expected sha %s but got %s",
			sh,
			sk.GetObjekteSha(),
		)

		err = errors.Wrapf(err, "Format: %T", format)

		return
	}

	return
}

func CalculateAndSetSha(
	sk SkuLikePtr,
	format objekte_format.Format,
) (err error) {
	w := sha.MakeWriter(nil)

	if _, err = format.FormatPersistentMetadatei(w, sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh := w.GetShaLike()
	sk.SetObjekteSha(sh)

	return
}

func ReadFromSha(
	sk SkuLikePtr,
	orf schnittstellen.ObjekteReaderFactory,
	format objekte_format.Format,
) (err error) {
	expected := sk.GetObjekteSha()

	var or schnittstellen.ShaReadCloser

	if or, err = orf.ObjekteReader(expected); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, or)

	if _, err = format.ParsePersistentMetadatei(or, sk); err != nil {
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
