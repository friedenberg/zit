package objekte_format

import (
	"bufio"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/echo/format"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
)

type v0 struct{}

func (f v0) FormatPersistentMetadatei(
	w1 io.Writer,
	c FormatterContext,
	o Options,
) (n int64, err error) {
	m := c.GetMetadatei()
	w := format.NewLineWriter()

	if o.Tai {
		w.WriteFormat("Tai %s", m.Tai)
	}

	w.WriteFormat("%s %s", gattung.Akte, &m.Akte)
	w.WriteFormat("%s %s", gattung.Typ, m.GetTyp())
	w.WriteFormat("%s %s", gattung.Bezeichnung, m.Bezeichnung)

	for _, e := range iter.SortedValues[kennung.Etikett](m.GetEtiketten()) {
		w.WriteFormat("%s %s", gattung.Etikett, e)
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
	_ Options,
) (n int64, err error) {
	m := c.GetMetadatei()

	etiketten := kennung.MakeEtikettMutableSet()

	r := bufio.NewReader(r1)

	typLineReader := ohio.MakeLineReaderIgnoreErrors(m.Typ.Set)

	esa := iter.MakeFuncSetString[kennung.Etikett, *kennung.Etikett](
		etiketten,
	)

	var g gattung.Gattung

	lr := format.MakeLineReaderConsumeEmpty(
		ohio.MakeLineReaderIterate(
			g.Set,
			ohio.MakeLineReaderKeyValues(
				map[string]interfaces.FuncSetString{
					"Tai":                        m.Tai.Set,
					gattung.Akte.String():        m.Akte.Set,
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

	m.SetEtiketten(etiketten)

	return
}
