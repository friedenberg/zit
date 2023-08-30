package objekte_format

import (
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/ohio"
	"github.com/friedenberg/zit/src/echo/format"
	"github.com/friedenberg/zit/src/echo/kennung"
)

type v3 struct{ includeTai bool }

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
		for _, e := range iter.SortedValues[kennung.Etikett](m.Etiketten) {
			w.WriteKeySpaceValue(gattung.Etikett, e)
		}
	}

	w.WriteKeySpaceValue("Gattung", c.GetKennungLike().GetGattung())
	w.WriteKeySpaceValue("Kennung", c.GetKennungLike())

	if f.includeTai {
		w.WriteKeySpaceValue("Tai", m.Tai)
	}

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

	var (
		g gattung.Gattung
		k kennung.Kennung
	)

	dr := ohio.MakeDelimReader('\n', r1)
	defer ohio.PutDelimReader(dr)

	var (
		lastKey string
		key     string
		val     string
	)

	for {
		key, val, err = dr.ReadOneKeyValue(" ")

		if err != nil {
			if errors.IsEOF(err) {
				err = nil
				break
			} else {
				err = errors.Wrap(err)
				return
			}
		}

		if key == "" {
			err = errors.Errorf("empty key at line %d", dr.Segments())
			return
		}

		if lastKey != "" && lastKey > key {
			err = errors.Errorf("keys not sorted")
			return
		}

		switch key {
		case "Akte":
			if err = m.AkteSha.Set(val); err != nil {
				err = errors.Wrap(err)
				return
			}

		case "Bezeichnung":
			if err = m.Bezeichnung.Set(val); err != nil {
				err = errors.Wrap(err)
				return
			}

		case "Etikett":
			if err = iter.AddString[kennung.Etikett, *kennung.Etikett](
				etiketten,
				val,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

		case "Gattung":
			if err = g.Set(val); err != nil {
				err = errors.Wrap(err)
				return
			}

		case "Kennung":
			if k, err = kennung.MakeWithGattung(g, val); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = c.SetKennungLike(k); err != nil {
				err = errors.Wrap(err)
				return
			}

		case "Tai":
			if err = m.Tai.Set(val); err != nil {
				err = errors.Wrap(err)
				return
			}

		case "Typ":
			if err = m.Typ.Set(val); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		lastKey = key
	}

	n = dr.N()

	if n == 0 {
		err = io.EOF
		return
	}

	m.Etiketten = etiketten

	c.SetMetadatei(m)

	return
}
