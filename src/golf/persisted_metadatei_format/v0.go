package persisted_metadatei_format

import (
	"bufio"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type v0 struct{}

func (f v0) FormatPersistentMetadatei(
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

func (f v0) ParsePersistentMetadatei(
	r1 io.Reader,
	c ParserContext,
) (n int64, err error) {
	m := c.GetMetadatei()

	etiketten := kennung.MakeEtikettMutableSet()

	r := bufio.NewReader(r1)

	typLineReader := format.MakeLineReaderIgnoreErrors(m.Typ.Set)

	esa := collections.MakeFuncSetString[kennung.Etikett, *kennung.Etikett](
		etiketten,
	)

	var g gattung.Gattung

	lr := format.MakeLineReaderConsumeEmpty(
		format.MakeLineReaderIterate(
			g.Set,
			format.MakeLineReaderKeyValues(
				map[string]schnittstellen.FuncSetString{
					"Tai":                        m.Tai.Set,
					gattung.Akte.String():        m.AkteSha.Set,
					gattung.Typ.String():         typLineReader,
					gattung.AkteTyp.String():     typLineReader,
					gattung.Bezeichnung.String(): m.Bezeichnung.Set,
					gattung.Etikett.String():     esa,
				},
			),
		),
	)

	if n, err = lr.ReadFrom(r); err != nil {
		err = errors.Wrap(err)
		return
	}

	m.Etiketten = etiketten.ImmutableClone()

	c.SetMetadatei(m)

	return
}
