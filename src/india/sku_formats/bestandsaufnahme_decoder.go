package sku_formats

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/log"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/ohio"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type FormatBestandsaufnahmeDecoder interface {
	ScanOne() (sku.SkuLikePtr, int64, error)
}

func MakeFormatbestandsaufnahmeDecoder(
	in io.Reader,
	of objekte_format.Format,
) FormatBestandsaufnahmeDecoder {
	return &bestandsaufnahmeDecoder{
		br:     ohio.MakeBoundaryReader(in, metadatei.Boundary+"\n"),
		format: of,
		es:     kennung.MakeEtikettMutableSet(),
	}
}

type bestandsaufnahmeDecoder struct {
	br         ohio.BoundaryReader
	format     objekte_format.Format
	afterFirst bool

	m  metadatei.Metadatei
	g  gattung.Gattung
	es kennung.EtikettMutableSet
	k  string
}

func (f *bestandsaufnahmeDecoder) ScanOne() (sk sku.SkuLikePtr, n int64, err error) {
	var (
		n1 int64
		n2 int
	)

	if !f.afterFirst {
		n2, err = f.br.ReadBoundary()
		n += int64(n2)

		if err != nil {
			if !errors.IsEOF(err) {
				err = errors.Wrap(err)
			}

			return
		}

		f.afterFirst = true
	}

	var h sku.Holder

	n1, err = f.format.ParsePersistentMetadatei(f.br, &h)
	n += n1

	if err != nil {
		if !errors.IsEOF(err) {
			err = errors.Wrap(err)
		}

		return
	}

	if sk, err = sku.MakeSkuLikeSansObjekteSha(
		h.Metadatei,
		h.KennungLike,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if sku.CalculateAndSetSha(sk, f.format); err != nil {
		err = errors.Wrap(err)
		return
	}

	n2, err = f.br.ReadBoundary()
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	log.Log().Printf("next boundary")

	return
}
