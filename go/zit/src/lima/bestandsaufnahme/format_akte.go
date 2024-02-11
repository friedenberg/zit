package bestandsaufnahme

import (
	"bufio"
	"io"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/golf/objekte_format"
	"code.linenisgreat.com/zit/src/india/sku_fmt"
)

type formatAkte struct {
	objekteFormat objekte_format.Format
	options       objekte_format.Options
}

func MakeAkteFormat(
	sv schnittstellen.StoreVersion,
	op objekte_format.Options,
) formatAkte {
	return formatAkte{
		objekteFormat: objekte_format.FormatForVersion(sv),
		options:       op,
	}
}

func (f formatAkte) ParseAkte(
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

func (f formatAkte) Format(w io.Writer, o *Akte) (n int64, err error) {
	return f.FormatParsedAkte(w, o)
}

func (f formatAkte) FormatParsedAkte(
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

	defer o.Skus.Restore()

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
