package store_util

import (
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/hotel/sku"
)

func (s *common) ReadOneExternalAkte(
	e sku.SkuLikeExternalPtr,
	t sku.SkuLikePtr,
) (err error) {
	e.SetMetadatei(t.GetMetadatei())

	var aw sha.WriteCloser

	if aw, err = s.AkteWriter(); err != nil {
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
