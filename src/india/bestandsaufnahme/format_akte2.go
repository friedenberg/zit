package bestandsaufnahme

import (
	"bufio"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/hotel/objekte"
)

type formatAkte2 struct {
	af            schnittstellen.AkteIOFactory
	objekteFormat objekte_format.Format
}

func (f formatAkte2) ParseSaveAkte(
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

	afterFirst := false
	var m metadatei.Metadatei
	var g gattung.Gattung
	es := kennung.MakeEtikettMutableSet()
	var k string

	if n, err = format.ReadLines(
		r,
		func(v string) (err error) {
			if v == metadatei.Boundary && afterFirst {
				var kl kennung.Kennung

				if kl, err = kennung.MakeWithGattung(g, k); err != nil {
					err = errors.Wrap(err)
					return
				}

				var sk sku.SkuLikePtr

				m.Etiketten = es.ImmutableClone()

				var m1 metadatei.Metadatei

				m1.ResetWith(m)

				if sk, err = sku.MakeSkuLikeSansObjekteSha(m1, kl); err != nil {
					err = errors.Wrap(err)
					return
				}

				if sku.CalculateAndSetSha(sk, f.objekteFormat); err != nil {
					err = errors.Wrap(err)
					return
				}

				if err = sku.AddSkuToHeap(&o.Skus, sk); err != nil {
					err = errors.Wrapf(err, "Sku: %s", sk)
					return
				}

				es.Reset()
				m.Reset()

				return
			} else if v == metadatei.Boundary {
				afterFirst = true
				return
			}

			idxSpace := strings.Index(v, " ")

			if idxSpace == -1 {
				err = errors.Errorf("expected to find space in line: %q", v)
				return
			}

			head := v[:idxSpace]
			tail := v[idxSpace+1:]

			switch head {
			case "Akte":
				return m.AkteSha.Set(tail)

			case "Bezeichnung":
				return m.Bezeichnung.Set(tail)

			case "Etikett":
				return collections.AddString[kennung.Etikett](es, tail)

			case "Gattung":
				return g.Set(tail)

			case "Kennung":
				k = tail
				return

			case "Tai":
				return m.Tai.Set(tail)

			case "Typ":
				return m.Typ.Set(tail)

			default:
				err = errors.Errorf("unsupported head %q for tail %q", head, tail)
				return
			}
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = aw.GetShaLike()

	return
}

func (f formatAkte2) Format(w io.Writer, o Akte) (n int64, err error) {
	return f.FormatParsedAkte(w, o)
}

func (f formatAkte2) FormatParsedAkte(
	w io.Writer,
	o Akte,
) (n int64, err error) {
	bw := bufio.NewWriter(w)
	defer errors.DeferredFlusher(&err, bw)

	fo := objekte.MakeFormatBestandsaufnahme(
		bw,
		objekte_format.BestandsaufnahmeFormatIncludeTai(),
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

		n1, err = fo.PrintOne(sk.SkuLikePtr)

		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
