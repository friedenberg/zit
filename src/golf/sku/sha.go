package sku

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/golf/objekte_format"
)

func CalculateAndSetSha(
	sk SkuLikePtr,
	format objekte_format.Format,
) (err error) {
	w := sha.MakeNopWriter()

	if _, err = format.FormatPersistentMetadatei(w, sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh := w.GetShaLike()
	sk.SetObjekteSha(sh)

	return
}
