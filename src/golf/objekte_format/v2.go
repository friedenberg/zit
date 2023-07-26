package objekte_format

import (
	"bufio"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/ohio"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type v2 struct{}

func (f v2) FormatPersistentMetadatei(
	w1 io.Writer,
	c FormatterContext,
) (n int64, err error) {
	m := c.GetMetadatei()
	w := format.NewLineWriter()

	if fcit, ok := c.(FormatterContextIncludeTai); ok && fcit.IncludeTai() {
		w.WriteFormat("Tai %s", m.Tai)
	}

	w.WriteFormat("%s %s", gattung.Akte, m.AkteSha)
	w.WriteFormat("%s %s", gattung.Typ, m.GetTyp())
	w.WriteFormat("%s %s", gattung.Bezeichnung, m.Bezeichnung)

	if m.Etiketten != nil {
		for _, e := range iter.SortedValues[kennung.Etikett](m.Etiketten) {
			w.WriteFormat("%s %s", gattung.Etikett, e)
		}
	}

	if n, err = w.WriteTo(w1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f v2) ParsePersistentMetadatei(
	r1 io.Reader,
	c ParserContext,
) (n int64, err error) {
	m := c.GetMetadatei()

	etiketten := kennung.MakeEtikettMutableSet()

	r := bufio.NewReader(r1)

	lr := format.MakeLineReaderConsumeEmpty(
		ohio.MakeLineReaderIterate(
			ohio.MakeLineReaderKeyValue("Tai", m.Tai.Set),
			ohio.MakeLineReaderKeyValue(gattung.Akte.String(), m.AkteSha.Set),
			ohio.MakeLineReaderKeyValue(
				gattung.Typ.String(),
				ohio.MakeLineReaderIgnoreErrors(m.Typ.Set),
			),
			ohio.MakeLineReaderKeyValue(
				gattung.Bezeichnung.String(),
				m.Bezeichnung.Set,
			),
			ohio.MakeLineReaderKeyValue(
				gattung.Etikett.String(),
				collections.MakeFuncSetString[kennung.Etikett, *kennung.Etikett](
					etiketten,
				),
			),
		),
	)

	if n, err = lr.ReadFrom(r); err != nil {
		err = errors.Wrap(err)
		return
	}

	m.Etiketten = etiketten

	c.SetMetadatei(m)

	return
}
