package sku

import (
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/ohio_ring_buffer2"
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
		var sb strings.Builder

		errors.PanicIfError(calculateAndSetSha(sk, format, o, &sb))

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
	return calculateAndSetSha(sk, format, o, nil)
}

func calculateAndSetSha(
	sk SkuLike,
	format objekte_format.Formatter,
	o objekte_format.Options,
	w1 io.Writer,
) (err error) {
	w := sha.MakeWriter(w1)

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

	if _, err = format.ParsePersistentMetadatei(
		ohio_ring_buffer2.MakeRingBuffer(or, 0),
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
