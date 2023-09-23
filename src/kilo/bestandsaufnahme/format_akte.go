package bestandsaufnahme

import (
	"bufio"
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/echo/format"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/sku_fmt"
)

type formatAkte struct {
	orfg                      schnittstellen.ObjekteReaderFactoryGetter
	persistentMetadateiFormat objekte_format.Format
	options                   objekte_format.Options
	af                        schnittstellen.AkteIOFactory
}

func (f formatAkte) ParseAkte(
	r io.Reader,
	o *Akte,
) (n int64, err error) {
	tml := sku.TryMakeSkuWithFormats(
		sku.MakeSkuFromLineTaiFirst,
		sku.MakeSkuFromLineGattungFirst,
	)

	if n, err = format.ReadLines(
		r,
		func(v string) (err error) {
			var sk *sku.Transacted2

			if sk, err = tml(v); err != nil {
				err = errors.Wrap(err)
				return
			}

			orf := f.orfg.ObjekteReaderFactory(sk)

			if err = sku.ReadFromSha(sk, orf, f.persistentMetadateiFormat, f.options); err != nil {
				err = errors.Wrap(err)
				return
			}

			return o.Skus.Add(sk)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f formatAkte) Format(w io.Writer, o Akte) (n int64, err error) {
	return f.FormatParsedAkte(w, o)
}

func (f formatAkte) FormatParsedAkte(w io.Writer, o Akte) (n int64, err error) {
	bw := bufio.NewWriter(w)
	defer errors.DeferredFlusher(&err, bw)

	defer func() {
		o.Skus.Restore()
	}()

	var n1 int

	for {
		sk, ok := o.Skus.PopAndSave()

		if !ok {
			break
		}

		l := fmt.Sprintf(
			"%s\n",
			sku_fmt.String(sk),
		)

		if n1, err = bw.WriteString(l); err != nil {
			err = errors.Wrap(err)
			return
		}

		n += int64(n1)
	}

	return
}
