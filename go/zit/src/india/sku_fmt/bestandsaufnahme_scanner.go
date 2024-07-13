package sku_fmt

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/go/zit/src/golf/ennui"
	"code.linenisgreat.com/zit/go/zit/src/golf/objekte_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type FormatBestandsaufnahmeScanner interface {
	Error() error
	GetTransacted() *sku.Transacted
	GetRange() ennui.Range
	Scan() bool
	SetDebug()
}

func MakeFormatBestandsaufnahmeScanner(
	in io.Reader,
	of objekte_format.Format,
	op objekte_format.Options,
) FormatBestandsaufnahmeScanner {
	return &bestandsaufnahmeScanner{
		ringBuffer: catgut.MakeRingBuffer(in, 0),
		format:     of,
		options:    op,
		es:         kennung.MakeTagMutableSet(),
	}
}

type bestandsaufnahmeScanner struct {
	ennui.Range

	ringBuffer *catgut.RingBuffer
	format     objekte_format.Format
	options    objekte_format.Options
	afterFirst bool

	m  metadatei.Metadatei
	g  gattung.Gattung
	es kennung.TagMutableSet
	k  string

	err     error
	lastSku *sku.Transacted
	debug   bool
}

func (scanner *bestandsaufnahmeScanner) SetDebug() {
	scanner.debug = true
}

func (scanner *bestandsaufnahmeScanner) Error() error {
	if errors.IsEOF(scanner.err) {
		return nil
	}

	return scanner.err
}

func (scanner *bestandsaufnahmeScanner) GetTransacted() *sku.Transacted {
	return scanner.lastSku
}

func (scanner *bestandsaufnahmeScanner) GetRange() ennui.Range {
	return scanner.Range
}

func (scanner *bestandsaufnahmeScanner) Scan() (ok bool) {
	if scanner.err != nil {
		return
	}

	var n1 int64

	scanner.lastSku = nil

	if !scanner.afterFirst {
		_, scanner.err = metadatei.ReadBoundary(scanner.ringBuffer)

		if errors.IsEOF(scanner.err) {
			return
		} else if scanner.err != nil {
			scanner.err = errors.Wrap(scanner.err)
			return
		}

		scanner.afterFirst = true
	}

	scanner.Offset += int64(len(metadatei.Boundary) + 1)
	scanner.ContentLength = 0

	scanner.lastSku = sku.GetTransactedPool().Get()

	scanner.ContentLength, scanner.err = scanner.format.ParsePersistentMetadatei(
		scanner.ringBuffer,
		scanner.lastSku,
		scanner.options,
	)

	if scanner.ContentLength == 0 {
		if scanner.err == io.EOF {
			return
		} else if scanner.err != nil {
			scanner.err = errors.Wrapf(scanner.err, "Bytes: %d", n1)
			scanner.err = errors.Wrapf(scanner.err, "Holder: %v", scanner.lastSku)
			return
		}
	}

	oldErr := scanner.err

	_, scanner.err = metadatei.ReadBoundary(scanner.ringBuffer)

	if errors.IsNotNilAndNotEOF(scanner.err) {
		scanner.err = errors.Wrap(errors.MakeMulti(scanner.err, oldErr))
		return
	} else if scanner.err == io.EOF {
		scanner.err = nil
		return
	}

	ok = true

	return
}
