package metadatei

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

type PersistedFormat struct {
	IgnoreTypErrors   bool
	EnforceFieldOrder bool
}

func (f PersistedFormat) Format(
	w1 io.Writer,
	c PersistentFormatterContext,
) (n int64, err error) {
	m := c.GetMetadatei()
	w := format.NewLineWriter()

	w.WriteFormat("%s %s", gattung.Akte, m.AkteSha)
	w.WriteFormat("%s %s", gattung.Typ, m.GetTyp())
	w.WriteFormat("%s %s", gattung.Bezeichnung, m.Bezeichnung)

	for _, e := range collections.SortedValues(m.Etiketten) {
		w.WriteFormat("%s %s", gattung.Etikett, e)
	}

	if n, err = w.WriteTo(w1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *PersistedFormat) Parse(
	r1 io.Reader,
	c PersistentParserContext,
) (n int64, err error) {
	m := c.GetMetadatei()

	etiketten := kennung.MakeEtikettMutableSet()

	r := bufio.NewReader(r1)

	typLineReader := m.Typ.Set

	if f.IgnoreTypErrors {
		typLineReader = format.MakeLineReaderIgnoreErrors(typLineReader)
	}

	esa := collections.MakeFuncSetString[kennung.Etikett, *kennung.Etikett](
		etiketten,
	)

	var g gattung.Gattung

	lineReaders := format.MakeLineReaderIterate(
		g.Set,
		format.MakeLineReaderKeyValues(
			map[string]schnittstellen.FuncSetString{
				gattung.Akte.String():        m.AkteSha.Set,
				gattung.Typ.String():         typLineReader,
				gattung.AkteTyp.String():     typLineReader,
				gattung.Bezeichnung.String(): m.Bezeichnung.Set,
				gattung.Etikett.String():     esa,
			},
		),
	)

	if f.EnforceFieldOrder {
		lineReaders = format.MakeLineReaderIterateStrict(
			g.Set, // will this work?
			format.MakeLineReaderKeyValue(gattung.Akte.String(), m.AkteSha.Set),
			format.MakeLineReaderKeyValue(gattung.Typ.String(), typLineReader),
			format.MakeLineReaderKeyValue(gattung.Bezeichnung.String(), m.Bezeichnung.Set),
			format.MakeLineReaderKeyValue(gattung.Etikett.String(), esa),
		)
	}

	lr := format.MakeLineReaderConsumeEmpty(lineReaders)

	if n, err = lr.ReadFrom(r); err != nil {
		err = errors.Wrap(err)
		return
	}

	m.Etiketten = etiketten.ImmutableClone()

	c.SetMetadatei(m)

	return
}
