package sku

import (
	"strings"

	"code.linenisgreat.com/zit-go/src/alfa/errors"
	"code.linenisgreat.com/zit-go/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit-go/src/charlie/catgut"
	"code.linenisgreat.com/zit-go/src/charlie/sha"
	"code.linenisgreat.com/zit-go/src/golf/objekte_format"
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
		var sb strings.Builder

		errors.PanicIfError(calculateAndSetSha(sk))

		err = errors.Errorf(
			"expected sha %q but got %q: used %q",
			sh,
			sk.GetObjekteSha(),
			sb.String(),
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
	return calculateAndSetSha(sk)
}

func calculateAndSetSha(
	sk SkuLike,
) (err error) {
	var actual *sha.Sha

	if actual, err = objekte_format.GetShaForMetadatei(
		objekte_format.Formats.Metadatei(),
		sk.GetMetadatei(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer sha.GetPool().Put(actual)

	if err = sk.SetObjekteSha(actual); err != nil {
		err = errors.Wrap(err)
		return
	}

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

	if _, err = format.ParsePersistentMetadatei(
		catgut.MakeRingBuffer(or, 0),
		sk,
		o,
	); err != nil {
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
