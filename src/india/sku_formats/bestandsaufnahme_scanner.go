package sku_formats

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/ohio"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type FormatBestandsaufnahmeScanner interface {
	Scan() (sku.SkuLikePtr, int64, error)
}

func MakeFormatbestandsaufnahmeScanner(
	in io.Reader,
	of objekte_format.Format,
) FormatBestandsaufnahmeScanner {
	return &bestandsaufnahmeScanner{
		br:     ohio.MakeBoundaryReader(in, metadatei.Boundary+"\n"),
		format: of,
		es:     kennung.MakeEtikettMutableSet(),
	}
}

type bestandsaufnahmeScanner struct {
	br         ohio.BoundaryReader
	format     objekte_format.Format
	afterFirst bool

	m  metadatei.Metadatei
	g  gattung.Gattung
	es kennung.EtikettMutableSet
	k  string
}

func (f *bestandsaufnahmeScanner) Scan() (sk sku.SkuLikePtr, n int64, err error) {
	var (
		n1 int64
		n2 int
	)

	if !f.afterFirst {
		n2, err = f.br.ReadBoundary()
		n += int64(n2)

		if errors.IsEOF(err) {
			err = nil
		} else if err != nil {
			err = errors.Wrap(err)
			return
		}

		f.afterFirst = true
	}

	var h sku.Holder

	n1, err = f.format.ParsePersistentMetadatei(f.br, &h)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if sk, err = sku.MakeSkuLikeSansObjekteSha(
		h.Metadatei,
		h.KennungLike,
	); err != nil {
		err = errors.Wrapf(err, "Bytes: %d", n1)
		err = errors.Wrapf(err, "Sku: %v", h)
		return
	}

	if sku.CalculateAndSetSha(sk, f.format); err != nil {
		err = errors.Wrap(err)
		return
	}

	n2, err = f.br.ReadBoundary()
	n += int64(n2)

	if errors.IsEOF(err) {
		err = io.EOF
		return
	} else if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
