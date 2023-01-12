package bestandsaufnahme

import (
	"bufio"
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/format"
)

type formatAkte struct {
}

func MakeFormatAkte() *formatAkte {
	return &formatAkte{}
}

func (f *formatAkte) Parse(r io.Reader, o *Objekte) (n int64, err error) {
	if n, err = format.ReadLines(
		r,
		o.Akte.Skus.AddString,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *formatAkte) Format(w io.Writer, o *Objekte) (n int64, err error) {
	sorted := o.Akte.Skus.SortedString()
	bw := bufio.NewWriter(w)
	defer errors.DeferredFlusher(&err, bw)

	for _, sk := range sorted {
		l := fmt.Sprintf(
			"%s\n",
			sk,
		)

		var n1 int

		if n1, err = bw.WriteString(l); err != nil {
			err = errors.Wrap(err)
			return
		}

		n += int64(n1)
	}

	return
}
