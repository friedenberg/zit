package bestandsaufnahme

import (
	"bufio"
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/foxtrot/sku"
)

type formatAkte struct {
	af schnittstellen.AkteIOFactory
}

func (f formatAkte) ParseSaveAkte(
	r1 io.Reader,
	o *Akte,
) (sh schnittstellen.ShaLike, n int64, err error) {
	var aw sha.WriteCloser

	if aw, err = f.af.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	r := io.TeeReader(r1, aw)

	tml := sku.TryMakeSkuWithFormats(
		sku.MakeSkuFromLineTaiFirst,
		sku.MakeSkuFromLineGattungFirst,
	)

	if n, err = format.ReadLines(
		r,
		func(v string) (err error) {
			var sk sku.SkuLike

			if sk, err = tml(v); err != nil {
				err = errors.Wrap(err)
				return
			}

			return sku.AddSkuToHeap(&o.Skus, sk)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = aw.GetShaLike()

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
			sku.String(sk),
		)

		if n1, err = bw.WriteString(l); err != nil {
			err = errors.Wrap(err)
			return
		}

		n += int64(n1)
	}

	return
}
