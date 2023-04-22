package persisted_metadatei_format

import (
	"bufio"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type v2 struct{}

func (f v2) Format(
	w1 io.Writer,
	c FormatterContext,
) (n int64, err error) {
	m := c.GetMetadatei()
	w := format.NewLineWriter()

	if !m.Tai.IsZero() {
		w.WriteFormat("Tai %s", m.Tai)
	}

	w.WriteFormat("%s %s", gattung.Akte, m.AkteSha)
	w.WriteFormat("%s %s", gattung.Typ, m.GetTyp())
	w.WriteFormat("%s %s", gattung.Bezeichnung, m.Bezeichnung)

	if m.Etiketten != nil {
		for _, e := range collections.SortedValues(m.Etiketten) {
			w.WriteFormat("%s %s", gattung.Etikett, e)
		}
	}

	if n, err = w.WriteTo(w1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f v2) Parse(
	r1 io.Reader,
	c ParserContext,
) (n int64, err error) {
	m := c.GetMetadatei()

	etiketten := kennung.MakeEtikettMutableSet()

	r := bufio.NewReader(r1)

	lr := format.MakeLineReaderConsumeEmpty(
		format.MakeLineReaderIterate(
			format.MakeLineReaderKeyValue("Tai", m.Tai.Set),
			format.MakeLineReaderKeyValue(gattung.Akte.String(), m.AkteSha.Set),
			format.MakeLineReaderKeyValue(
				gattung.Typ.String(),
				format.MakeLineReaderIgnoreErrors(m.Typ.Set),
			),
			format.MakeLineReaderKeyValue(gattung.Bezeichnung.String(), m.Bezeichnung.Set),
			format.MakeLineReaderKeyValue(
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
