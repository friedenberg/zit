package zettel

import (
	"bufio"
	"crypto/sha256"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/echo/format"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/schnittstellen"
)

func (z Objekte) ObjekteSha() (s sha.Sha, err error) {
	errors.TodoP1("remove this")

	hash := sha256.New()

	o := FormatObjekte{}

	c := ObjekteFormatterContext{
		Zettel: z,
	}

	if _, err = o.Format(hash, &c.Zettel); err != nil {
		err = errors.Wrap(err)
		return
	}

	s = sha.FromHash(hash)

	return
}

type FormatObjekte struct {
	IgnoreTypErrors bool
}

func (f FormatObjekte) Format(
	w1 io.Writer,
	z *Objekte,
) (n int64, err error) {
	errors.TodoP1("replace with objekte.Format")

	w := format.NewLineWriter()

	w.WriteFormat("%s %s", gattung.Akte, z.Akte)
	w.WriteFormat("%s %s", gattung.Typ, z.Typ)
	w.WriteFormat("%s %s", gattung.Bezeichnung, z.Bezeichnung)

	for _, e := range z.Etiketten.Sorted() {
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

	typLineReader := z.Typ.Set

	if f.IgnoreTypErrors {
		typLineReader = format.MakeLineReaderIgnoreErrors(typLineReader)
	}

	esa := collections.MakeFuncSetString[kennung.Etikett, *kennung.Etikett](etiketten)

	if n, err = format.ReadLines(
		r,
		format.MakeLineReaderRepeat(
			format.MakeLineReaderKeyValues(
				map[string]schnittstellen.FuncSetString{
					gattung.Akte.String():        z.Akte.Set,
					gattung.Typ.String():         typLineReader,
					gattung.AkteTyp.String():     typLineReader,
					gattung.Bezeichnung.String(): z.Bezeichnung.Set,
					gattung.Etikett.String():     esa,
				},
			),
		),
		// format.MakeLineReaderKeyValue(gattung.Akte.String(), z.Akte.Set),
		// format.MakeLineReaderKeyValue(gattung.Typ.String(), typLineReader),
		// format.MakeLineReaderKeyValue(gattung.Bezeichnung.String(), z.Bezeichnung.Set),
		// format.MakeLineReaderKeyValue(gattung.Etikett.String(), etiketten.AddString),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	z.Etiketten = etiketten.Copy()

	return
}
