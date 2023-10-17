package to_merge

import (
	"io"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/ohio"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/golf/objekte_format"
)

type bestandsaufnahmePrinter struct {
	format            objekte_format.Formatter
	options           objekte_format.Options
	out               io.Writer
	firstBoundaryOnce *sync.Once
}

type FormatBestandsaufnahmePrinter interface {
	Print(objekte_format.FormatterContext) (int64, error)
	PrintMany(...objekte_format.FormatterContext) (int64, error)
}

func MakeFormatBestandsaufnahmePrinter(
	out io.Writer,
	of objekte_format.Formatter,
	op objekte_format.Options,
) FormatBestandsaufnahmePrinter {
	return bestandsaufnahmePrinter{
		format:            of,
		options:           op,
		out:               out,
		firstBoundaryOnce: &sync.Once{},
	}
}

func (f bestandsaufnahmePrinter) printBoundary() (n int64, err error) {
	if n, err = ohio.WriteLine(f.out, metadatei.Boundary); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f bestandsaufnahmePrinter) printFirstBoundary() (n int64, err error) {
	f.firstBoundaryOnce.Do(
		func() {
			n, err = f.printBoundary()
		},
	)

	return
}

func (f bestandsaufnahmePrinter) PrintMany(
	tlps ...objekte_format.FormatterContext,
) (n int64, err error) {
	for _, tlp := range tlps {
		var n1 int64
		n1, err = f.Print(tlp)
		n += n1

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (f bestandsaufnahmePrinter) Print(
	tlp objekte_format.FormatterContext,
) (n int64, err error) {
	pfs := [3]func() (int64, error){
		f.printFirstBoundary,
		func() (int64, error) {
			return f.format.FormatPersistentMetadatei(f.out, tlp, f.options)
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
