package zettel

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

type FormatObjekte struct {
	IgnoreTypErrors   bool
	EnforceFieldOrder bool
}

func (f FormatObjekte) Format(
	w1 io.Writer,
	z *Objekte,
) (n int64, err error) {
	errors.TodoP1("replace with objekte.Format")

	w := format.NewLineWriter()

	w.WriteFormat("%s %s", gattung.Akte, z.Akte)
	w.WriteFormat("%s %s", gattung.Typ, z.GetTyp())
	w.WriteFormat("%s %s", gattung.Bezeichnung, z.Metadatei.Bezeichnung)

	for _, e := range collections.SortedValues(z.Metadatei.Etiketten) {
		w.WriteFormat("%s %s", gattung.Etikett, e)
	}

	if n, err = w.WriteTo(w1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *FormatObjekte) Parse(
	r1 io.Reader,
	z *Objekte,
) (n int64, err error) {
	etiketten := kennung.MakeEtikettMutableSet()

	r := bufio.NewReader(r1)

	typLineReader := z.Metadatei.Typ.Set

	if f.IgnoreTypErrors {
		typLineReader = format.MakeLineReaderIgnoreErrors(typLineReader)
	}

	esa := collections.MakeFuncSetString[kennung.Etikett, *kennung.Etikett](
		etiketten,
	)

	lineReaders := []schnittstellen.FuncSetString{
		format.MakeLineReaderRepeat(
			format.MakeLineReaderKeyValues(
				map[string]schnittstellen.FuncSetString{
					gattung.Akte.String():        z.Akte.Set,
					gattung.Typ.String():         typLineReader,
					gattung.AkteTyp.String():     typLineReader,
					gattung.Bezeichnung.String(): z.Metadatei.Bezeichnung.Set,
					gattung.Etikett.String():     esa,
				},
			),
		),
	}

	if f.EnforceFieldOrder {
		lineReaders = []schnittstellen.FuncSetString{
			format.MakeLineReaderKeyValue(gattung.Akte.String(), z.Akte.Set),
			format.MakeLineReaderKeyValue(gattung.Typ.String(), typLineReader),
			format.MakeLineReaderKeyValue(gattung.Bezeichnung.String(), z.Metadatei.Bezeichnung.Set),
			format.MakeLineReaderKeyValue(gattung.Etikett.String(), esa),
		}
	}

	if n, err = format.ReadLines(r, lineReaders...); err != nil {
		err = errors.Wrap(err)
		return
	}

	z.Metadatei.Etiketten = etiketten.ImmutableClone()

	return
}
