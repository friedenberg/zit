package sku_fmt

import (
	"io"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
)

type FormatBestandsaufnahmePrinter interface {
	Offset() int64
	Print(object_inventory_format.FormatterContext) (int64, error)
	PrintMany(...object_inventory_format.FormatterContext) (int64, error)
}

type bestandsaufnahmePrinter struct {
	format            object_inventory_format.Formatter
	options           object_inventory_format.Options
	out               io.Writer
	offset            int64
	firstBoundaryOnce *sync.Once
}

func MakeFormatInventoryListPrinter(
	out io.Writer,
	of object_inventory_format.Formatter,
	op object_inventory_format.Options,
) FormatBestandsaufnahmePrinter {
	return &bestandsaufnahmePrinter{
		format:            of,
		options:           op,
		out:               out,
		offset:            int64(len(object_metadata.Boundary) + 1),
		firstBoundaryOnce: &sync.Once{},
	}
}

func (f *bestandsaufnahmePrinter) printBoundary() (n int64, err error) {
	if n, err = ohio.WriteLine(f.out, object_metadata.Boundary); err != nil {
		err = errors.Wrap(err)
		return
	}

	f.offset += n

	return
}

func (f *bestandsaufnahmePrinter) printFirstBoundary() (n int64, err error) {
	f.firstBoundaryOnce.Do(
		func() {
			if n, err = ohio.WriteLine(f.out, object_metadata.Boundary); err != nil {
				err = errors.Wrap(err)
				return
			}
		},
	)

	return
}

func (f *bestandsaufnahmePrinter) PrintMany(
	tlps ...object_inventory_format.FormatterContext,
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

func (f *bestandsaufnahmePrinter) Offset() int64 {
	return f.offset
}

func (f *bestandsaufnahmePrinter) makeFuncFormatOne(
	tlp object_inventory_format.FormatterContext,
) func() (int64, error) {
	return func() (int64, error) {
		n1, err := f.format.FormatPersistentMetadata(f.out, tlp, f.options)
		f.offset += n1
		return n1, err
	}
}

func (f *bestandsaufnahmePrinter) Print(
	tlp object_inventory_format.FormatterContext,
) (n int64, err error) {
	var n1 int64

	pfs := [3]func() (int64, error){
		f.printFirstBoundary,
		f.makeFuncFormatOne(tlp),
		f.printBoundary,
	}

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
