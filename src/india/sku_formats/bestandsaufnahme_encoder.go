package sku_formats

import (
	"io"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/ohio"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/golf/objekte_format"
)

type bestandsaufnahmeEncoder struct {
	format            objekte_format.Formatter
	out               io.Writer
	firstBoundaryOnce *sync.Once
}

type FormatBestandsaufnahmeEncoder interface {
	PrintOne(objekte_format.FormatterContext) (int64, error)
}

func MakeFormatBestandsaufnahmeEncoder(
	out io.Writer,
	of objekte_format.Formatter,
) FormatBestandsaufnahmeEncoder {
	return bestandsaufnahmeEncoder{
		format:            of,
		out:               out,
		firstBoundaryOnce: &sync.Once{},
	}
}

func (f bestandsaufnahmeEncoder) printBoundary() (n int64, err error) {
	if n, err = ohio.WriteLine(f.out, metadatei.Boundary); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f bestandsaufnahmeEncoder) printFirstBoundary() (n int64, err error) {
	f.firstBoundaryOnce.Do(
		func() {
			n, err = f.printBoundary()
		},
	)

	return
}

func (f bestandsaufnahmeEncoder) PrintOne(
	tlp objekte_format.FormatterContext,
) (n int64, err error) {
	pfs := [3]func() (int64, error){
		f.printFirstBoundary,
		func() (int64, error) {
			return f.format.FormatPersistentMetadatei(f.out, tlp)
		},
		f.printBoundary,
	}

	var n1 int64

	for _, pf := range pfs {
		n1, err = pf()

		n += n1

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
