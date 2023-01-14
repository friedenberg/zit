package bestandsaufnahme

import (
	"bufio"
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/golf/sku"
)

type formatAkte struct {
}

func MakeFormatAkte() *formatAkte {
	return &formatAkte{}
}

func (f *formatAkte) Parse(r io.Reader, o *Objekte) (n int64, err error) {
	if n, err = format.ReadLines(
		r,
		func(v string) (err error) {
			return collections.AddString[sku.Sku2, *sku.Sku2](&o.Akte.Skus, v)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *formatAkte) Format(w io.Writer, o *Objekte) (n int64, err error) {
	bw := bufio.NewWriter(w)
	defer errors.DeferredFlusher(&err, bw)

	defer func() {
		o.Akte.Skus.Restore()
	}()

	var n1 int

	for {
		sk, ok := o.Akte.Skus.PopAndSave()

		if !ok {
			break
		}

		l := fmt.Sprintf(
			"%s\n",
			sk,
		)

		if n1, err = bw.WriteString(l); err != nil {
			err = errors.Wrap(err)
			return
		}

		n += int64(n1)
	}

	return
}
