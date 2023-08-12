package sku

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
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
