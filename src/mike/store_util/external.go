package store_util

import (
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type ExternalReader interface {
	ReadOneExternalObjekte(
		e sku.SkuLikeExternalPtr,
		t sku.SkuLikePtr,
	) (err error)

	ReadOneExternalAkte(
		e sku.SkuLikeExternalPtr,
		t sku.SkuLikePtr,
	) (err error)
}

func (s *common) ReadOneExternalObjekte(
	e sku.SkuLikeExternalPtr,
	t sku.SkuLikePtr,
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
