package sku

import (
	"strings"

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
		return
	}

	return
}

func CalculateAndSetSha(
	sk SkuLikePtr,
	format objekte_format.Format,
) (err error) {
	st := &strings.Builder{}
	w := sha.MakeWriter(st)

	if _, err = format.FormatPersistentMetadatei(w, sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh := w.GetShaLike()
	sk.SetObjekteSha(sh)

	return
}
