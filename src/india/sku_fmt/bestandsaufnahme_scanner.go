package sku_fmt

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
	Error() error
	GetSkuLikePtr() sku.SkuLikePtr
	Scan() bool
}

func MakeFormatBestandsaufnahmeScanner(
	in io.Reader,
	of objekte_format.Format,
	op objekte_format.Options,
) FormatBestandsaufnahmeScanner {
	return &bestandsaufnahmeScanner{
		br:      ohio.MakeBoundaryReader(in, metadatei.Boundary+"\n"),
		format:  of,
		options: op,
		es:      kennung.MakeEtikettMutableSet(),
	}
}

func MakeFormatBestandsaufnahmeScanner2(
	in io.Reader,
	of objekte_format.Format,
	op objekte_format.Options,
) FormatBestandsaufnahmeScanner {
	return &bestandsaufnahmeScanner{
		br:             ohio.MakeBoundaryReader(in, metadatei.Boundary+"\n"),
		format:         of,
		options:        op,
		es:             kennung.MakeEtikettMutableSet(),
		useTransacted2: true,
	}
}

type bestandsaufnahmeScanner struct {
	br         ohio.BoundaryReader
	format     objekte_format.Format
	options    objekte_format.Options
	afterFirst bool

	useTransacted2 bool

	m  metadatei.Metadatei
	g  gattung.Gattung
	es kennung.EtikettMutableSet
	k  string

	err     error
	lastSku sku.SkuLikePtr
	lastN   int64
}

func (f *bestandsaufnahmeScanner) Error() error {
	if errors.IsEOF(f.err) {
		return nil
	}

	return f.err
}

func (f *bestandsaufnahmeScanner) GetSkuLikePtr() sku.SkuLikePtr {
	return f.lastSku
}

func (f *bestandsaufnahmeScanner) Scan() (ok bool) {
	if f.err != nil {
		return
	}

	var (
		n1 int64
		n2 int
	)

	f.lastN = 0
	f.lastSku = nil

	if !f.afterFirst {
		n2, f.err = f.br.ReadBoundary()
		f.lastN += int64(n2)

		if errors.IsEOF(f.err) {
			return
		} else if f.err != nil {
			f.err = errors.Wrap(f.err)
			return
		}

		f.afterFirst = true
	}

	var h sku.Holder

	n1, f.err = f.format.ParsePersistentMetadatei(f.br, &h, f.options)
	f.lastN += n1

	if errors.IsEOF(f.err) {
		f.err = errors.Errorf("unexpected eof")
		return
	} else if f.err != nil {
		f.err = errors.Wrapf(f.err, "Bytes: %d", n1)
		f.err = errors.Wrapf(f.err, "Holder: %v", h)
		return
	}

	if f.useTransacted2 {
		if f.lastSku, f.err = sku.MakeSkuLikeSansObjekteSha2(
			h.Metadatei,
			h.KennungLike,
		); f.err != nil {
			f.err = errors.Wrapf(f.err, "Bytes: %d", n1)
			f.err = errors.Wrapf(f.err, "Sku: %v", h)
			return
		}
	} else {
		if f.lastSku, f.err = sku.MakeSkuLikeSansObjekteSha(
			h.Metadatei,
			h.KennungLike,
		); f.err != nil {
			f.err = errors.Wrapf(f.err, "Bytes: %d", n1)
			f.err = errors.Wrapf(f.err, "Sku: %v", h)
			return
		}
	}

	if f.err = sku.CalculateAndSetSha(f.lastSku, f.format, f.options.SansVerzeichnisse()); f.err != nil {
		f.err = errors.Wrap(f.err)
		return
	}

	n2, f.err = f.br.ReadBoundary()
	f.lastN += int64(n2)

	if f.err != nil && !errors.IsEOF(f.err) {
		f.err = errors.Wrap(f.err)
		return
	}

	ok = true

	return
}
