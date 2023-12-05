package sku_fmt

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/catgut"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type FormatBestandsaufnahmeScanner interface {
	Error() error
	GetTransacted() *sku.Transacted
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
		es:         kennung.MakeEtikettMutableSet(),
	}
}

type bestandsaufnahmeScanner struct {
	ringBuffer *catgut.RingBuffer
	format     objekte_format.Format
	options    objekte_format.Options
	afterFirst bool

	m  metadatei.Metadatei
	g  gattung.Gattung
	es kennung.EtikettMutableSet
	k  string

	err     error
	lastSku *sku.Transacted
	lastN   int64
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

func (scanner *bestandsaufnahmeScanner) Scan() (ok bool) {
	if scanner.err != nil {
		return
	}

	var (
		n1 int64
		n2 int
	)

	scanner.lastN = 0
	scanner.lastSku = nil

	if !scanner.afterFirst {
		n2, scanner.err = metadatei.ReadBoundary(scanner.ringBuffer)
		scanner.lastN += int64(n2)

		if errors.IsEOF(scanner.err) {
			return
		} else if scanner.err != nil {
			scanner.err = errors.Wrap(scanner.err)
			return
		}

		scanner.afterFirst = true
	}

	scanner.lastSku = sku.GetTransactedPool().Get()

	n1, scanner.err = scanner.format.ParsePersistentMetadatei(
		scanner.ringBuffer,
		scanner.lastSku,
		scanner.options,
	)

	scanner.lastN += n1

	if n1 == 0 {
		if scanner.err == io.EOF {
			return
		} else if scanner.err != nil {
			scanner.err = errors.Wrapf(scanner.err, "Bytes: %d", n1)
			scanner.err = errors.Wrapf(scanner.err, "Holder: %v", scanner.lastSku)
			return
		}
	}

	scanner.lastSku.SetObjekteSha(&scanner.lastSku.Metadatei.Verzeichnisse.Sha)

	oldErr := scanner.err

	n2, scanner.err = metadatei.ReadBoundary(scanner.ringBuffer)
	scanner.lastN += int64(n2)

	if scanner.err != nil && !errors.IsEOF(scanner.err) {
		scanner.err = errors.Wrap(errors.MakeMulti(scanner.err, oldErr))
		return
	} else if scanner.err == io.EOF {
		scanner.err = nil
		return
	}

	ok = true

	return
}
