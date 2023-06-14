package objekte_format

import (
	"bufio"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/ohio"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type v3 struct{}

func (f v3) FormatPersistentMetadatei(
	w1 io.Writer,
	c FormatterContext,
) (n int64, err error) {
	m := c.GetMetadatei()

	w := format.NewLineWriter()

	if !m.AkteSha.IsNull() {
		w.WriteKeySpaceValue(gattung.Akte, m.AkteSha)
	}

	lines := strings.Split(m.Bezeichnung.String(), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}
		w.WriteKeySpaceValue(gattung.Bezeichnung, line)
	}

	if m.Etiketten != nil {
		for _, e := range collections.SortedValues(m.Etiketten) {
			w.WriteKeySpaceValue(gattung.Etikett, e)
		}
	}

	w.WriteKeySpaceValue("Gattung", c.GetKennung().GetGattung())
	w.WriteKeySpaceValue("Kennung", c.GetKennung())
	w.WriteKeySpaceValue("Tai", m.Tai)

	if !m.Typ.IsEmpty() {
		w.WriteKeySpaceValue(gattung.Typ, m.GetTyp())
	}

	if n, err = w.WriteTo(w1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO-P1 implement proper parsing of v3 format
func (f v3) ParsePersistentMetadatei(
	r1 io.Reader,
	c ParserContext,
) (n int64, err error) {
	m := c.GetMetadatei()

	etiketten := kennung.MakeEtikettMutableSet()

	r := bufio.NewReader(r1)

	lr := format.MakeLineReaderConsumeEmpty(
		ohio.MakeLineReaderIterate(
			ohio.MakeLineReaderKeyValue(gattung.Akte.String(), m.AkteSha.Set),
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
			func(v string) (err error) {
				var k kennung.Kennung

				if k, err = kennung.Make(v); err != nil {
					err = errors.Wrap(err)
					return
				}

				if err = c.SetKennung(k); err != nil {
					err = errors.Wrap(err)
					return
				}

				return
			},
			ohio.MakeLineReaderKeyValue("Tai", m.Tai.Set),
			ohio.MakeLineReaderKeyValue(
				gattung.Typ.String(),
				ohio.MakeLineReaderIgnoreErrors(m.Typ.Set),
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
