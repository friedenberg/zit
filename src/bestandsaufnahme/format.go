package bestandsaufnahme

import (
	"bufio"
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	format_pkg "github.com/friedenberg/zit/src/delta/format"
)

type format struct {
}

func MakeFormat() *format {
	return &format{}
}

func (f *format) Parse(r io.Reader, b *Akte) (n int64, err error) {
	if n, err = format_pkg.ReadLines(
		r,
		b.Skus.AddString,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *format) Format(w io.Writer, b *Akte) (n int64, err error) {
	sorted := b.Skus.SortedString()
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
