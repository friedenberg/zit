package bestandsaufnahme

import (
	"bufio"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/india/sku_fmt"
)

type formatAkte2 struct {
	objekteFormat objekte_format.Format
	options       objekte_format.Options
}

func MakeAkteFormat(
	sv schnittstellen.StoreVersion,
	op objekte_format.Options,
) formatAkte2 {
	return formatAkte2{
		objekteFormat: objekte_format.FormatForVersion(sv),
		options:       op,
	}
}

func (f formatAkte2) ParseAkte(
	r io.Reader,
	o *Akte,
) (n int64, err error) {
	dec := sku_fmt.MakeFormatBestandsaufnahmeScanner(
		r,
		f.objekteFormat,
		f.options,
	)

	// dec.SetDebug()

	for dec.Scan() {
		sk := dec.GetTransacted()

		if err = o.Skus.Add(sk); err != nil {
			err = errors.Wrapf(err, "Sku: %s", sk)
			return
		}
	}

	if err = dec.Error(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f formatAkte2) Format(w io.Writer, o *Akte) (n int64, err error) {
	return f.FormatParsedAkte(w, o)
}

func (f formatAkte2) FormatParsedAkte(
	w io.Writer,
	o *Akte,
) (n int64, err error) {
	bw := bufio.NewWriter(w)
	defer errors.DeferredFlusher(&err, bw)

	fo := sku_fmt.MakeFormatBestandsaufnahmePrinter(
		bw,
		f.objekteFormat,
		f.options,
	)

	defer func() {
		o.Skus.Restore()
	}()

	var n1 int64

	for {
		sk, ok := o.Skus.PopAndSave()

		if !ok {
			break
		}

		n1, err = fo.Print(sk)

		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
