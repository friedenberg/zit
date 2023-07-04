package objekte

import (
	"io"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/ohio"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/golf/objekte_format"
)

type formatBestandsaufnahme struct {
	format            objekte_format.Formatter
	out               io.Writer
	firstBoundaryOnce *sync.Once
}

type FormatBestandsaufnahme interface {
	PrintOne(TransactedLikePtr) (int64, error)
}

func MakeFormatBestandsaufnahme(
	out io.Writer,
	of objekte_format.Formatter,
) FormatBestandsaufnahme {
	return formatBestandsaufnahme{
		format:            of,
		out:               out,
		firstBoundaryOnce: &sync.Once{},
	}
}

func (f formatBestandsaufnahme) printBoundary() (n int64, err error) {
	if n, err = ohio.WriteLine(f.out, metadatei.Boundary); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f formatBestandsaufnahme) printFirstBoundary() (n int64, err error) {
	f.firstBoundaryOnce.Do(
		func() {
			n, err = f.printBoundary()
		},
	)

	return
}

func (f formatBestandsaufnahme) PrintOne(
	tlp TransactedLikePtr,
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
