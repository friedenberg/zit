package objekte_format

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/ohio"
	"github.com/friedenberg/zit/src/echo/kennung"
)

type v4 struct{}

func (f v4) FormatPersistentMetadatei(
	w1 io.Writer,
	c FormatterContext,
	o Options,
) (n int64, err error) {
	w := ohio.GetPoolBufioWriter().Get()
	defer ohio.GetPoolBufioWriter().Put(w)

	w.Reset(w1)
	defer errors.DeferredFlusher(&err, w)

	m := c.GetMetadatei()

	mh := sha.MakeWriter(nil)
	mw := io.MultiWriter(w, mh)
	var (
		n1 int
		n2 int64
	)

	if !m.AkteSha.IsNull() {
		n1, err = ohio.WriteKeySpaceValueNewlineString(
			mw,
			gattung.Akte.String(),
			m.AkteSha.String(),
		)
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	lines := strings.Split(m.Bezeichnung.String(), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		n1, err = ohio.WriteKeySpaceValueNewlineString(
			mw,
			gattung.Bezeichnung.String(),
			line,
		)
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	es := m.GetEtiketten()

	for _, e := range iter.SortedValues[kennung.Etikett](es) {
		n1, err = ohio.WriteKeySpaceValueNewlineString(
			mw,
			gattung.Etikett.String(),
			e.String(),
		)
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	n1, err = ohio.WriteKeySpaceValueNewlineString(
		w,
		"Gattung",
		c.GetKennungLike().GetGattung().GetGattungString(),
	)
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = ohio.WriteKeySpaceValueNewlineString(
		w,
		"Kennung",
		c.GetKennungLike().String(),
	)
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if o.IncludeTai {
		n1, err = ohio.WriteKeySpaceValueNewlineString(
			mw,
			"Tai",
			m.Tai.String(),
		)
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if !m.Typ.IsEmpty() {
		n1, err = ohio.WriteKeySpaceValueNewlineString(
			mw,
			gattung.Typ.String(),
			m.GetTyp().String(),
		)
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if o.IncludeVerzeichnisse {
		if m.Verzeichnisse.Archiviert.Bool() {
			n1, err = ohio.WriteKeySpaceValueNewlineString(
				w,
				"Verzeichnisse-Archiviert",
				m.Verzeichnisse.Archiviert.String(),
			)
			n += int64(n1)

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		if m.Verzeichnisse.GetExpandedEtiketten().Len() > 0 {
			k := fmt.Sprintf(
				"Verzeichnisse-%s-Expanded",
				gattung.Etikett.String(),
			)
			for _, e := range iter.SortedValues[kennung.Etikett](m.Verzeichnisse.GetExpandedEtiketten()) {
				n1, err = ohio.WriteKeySpaceValueNewlineString(
					w,
					k,
					e.String(),
				)
				n += int64(n1)

				if err != nil {
					err = errors.Wrap(err)
					return
				}
			}
		}

		if m.Verzeichnisse.GetImplicitEtiketten().Len() > 0 {
			k := fmt.Sprintf(
				"Verzeichnisse-%s-Implicit",
				gattung.Etikett.String(),
			)

			for _, e := range iter.SortedValues[kennung.Etikett](m.Verzeichnisse.GetImplicitEtiketten()) {
				n2, err = ohio.WriteKeySpaceValueNewline(
					w,
					k,
					e.Bytes(),
				)
				n += int64(n2)

				if err != nil {
					err = errors.Wrap(err)
					return
				}
			}
		}

		if !m.Verzeichnisse.Mutter.IsNull() {
			n1, err = ohio.WriteKeySpaceValueNewlineString(
				mw,
				"Verzeichnisse-Mutter",
				m.Verzeichnisse.Mutter.String(),
			)
			n += int64(n1)

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		actual := mh.GetShaLike()
		// TODO-P1 set value

		// if !m.Verzeichnisse.Sha.IsNull() &&
		// 	!m.Verzeichnisse.Sha.EqualsSha(actual) {
		// 	err = errors.Errorf(
		// 		"expected %q but got %q -> %q",
		// 		m.Verzeichnisse.Sha,
		// 		actual,
		// 		sb.String(),
		// 	)
		// 	return
		// }

		n1, err = ohio.WriteKeySpaceValueNewlineString(
			w,
			"Verzeichnisse-Sha",
			actual.String(),
		)
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (f v4) ParsePersistentMetadatei(
	r1 io.Reader,
	c ParserContext,
	o Options,
) (n int64, err error) {
	m := c.GetMetadatei()

	var (
		g gattung.Gattung
		k kennung.Kennung2
	)

	dr := ohio.MakeDelimReader('\n', r1)
	defer ohio.PutDelimReader(dr)

	var (
		lastKey []byte
		key     []byte
		val     []byte
	)

	mh := sha.MakeWriter(nil)

	for {
		key, val, err = dr.ReadOneKeyValueBytes(' ')

		if err != nil {
			if errors.IsEOF(err) {
				err = nil
				break
			} else {
				err = errors.Wrap(err)
				return
			}
		}

		if len(key) == 0 {
			err = errors.Errorf("empty key at line %d", dr.Segments())
			return
		}

		if len(lastKey) >= 0 && bytes.Compare(lastKey, key) == 1 {
			err = errors.Errorf("keys not sorted")
			return
		}

		writeMetadateiHashString := false

		// this cast in an iff costs nothing, vs a switch: https://github.com/golang/go/issues/24937
		if string(key) == "Akte" {
			if err = m.AkteSha.Set(string(val)); err != nil {
				err = errors.Wrap(err)
				return
			}

			writeMetadateiHashString = true

		} else if string(key) == "Bezeichnung" {
			if err = m.Bezeichnung.Set(string(val)); err != nil {
				err = errors.Wrap(err)
				return
			}

			writeMetadateiHashString = true

		} else if string(key) == "Etikett" {
			e := kennung.GetEtikettPool().Get()

			if err = e.Set(string(val)); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = m.AddEtikettPtr(e); err != nil {
				err = errors.Wrap(err)
				return
			}

			writeMetadateiHashString = true

		} else if string(key) == "Gattung" {
			if err = g.Set(string(val)); err != nil {
				err = errors.Wrap(err)
				return
			}

		} else if string(key) == "Kennung" {
			if err = k.SetWithGattung(string(val), g); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = c.SetKennungLike(k); err != nil {
				err = errors.Wrap(err)
				return
			}

		} else if string(key) == "Tai" {
			if err = m.Tai.Set(string(val)); err != nil {
				err = errors.Wrap(err)
				return
			}

			writeMetadateiHashString = true

		} else if string(key) == "Typ" {
			if err = m.Typ.Set(string(val)); err != nil {
				err = errors.Wrap(err)
				return
			}

			writeMetadateiHashString = true

		} else if string(key) == "Verzeichnisse-Archiviert" {
			if err = m.Verzeichnisse.Archiviert.Set(string(val)); err != nil {
				err = errors.Wrap(err)
				return
			}

		} else if string(key) == "Verzeichnisse-Etikett-Implicit" {
			if !o.IncludeVerzeichnisse {
				err = errors.Errorf(
					"format specifies not to include Verzeichnisse but found %q",
					key,
				)
				return
			}

			e := kennung.GetEtikettPool().Get()

			if err = e.Set(string(val)); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = m.Verzeichnisse.GetImplicitEtikettenMutable().AddPtr(e); err != nil {
				err = errors.Wrap(err)
				return
			}

		} else if string(key) == "Verzeichnisse-Etikett-Expanded" {
			if !o.IncludeVerzeichnisse {
				err = errors.Errorf(
					"format specifies not to include Verzeichnisse but found %q",
					key,
				)
				return
			}

			e := kennung.GetEtikettPool().Get()

			if err = e.Set(string(val)); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = m.Verzeichnisse.GetExpandedEtikettenMutable().AddPtr(e); err != nil {
				err = errors.Wrap(err)
				return
			}
		} else if string(key) == "Verzeichnisse-Mutter" {
			if err = m.Verzeichnisse.Mutter.Set(string(val)); err != nil {
				err = errors.Wrap(err)
				return
			}

			writeMetadateiHashString = true

		} else if string(key) == "Verzeichnisse-Sha" {
			if err = m.Verzeichnisse.Sha.Set(string(val)); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		lastKey = key

		if !writeMetadateiHashString {
			continue
		}

		if _, err = ohio.WriteKeySpaceValueNewline(mh, string(key), val); err != nil {
			err = errors.Wrap(err)
			return
		}

		writeMetadateiHashString = false
	}

	n = dr.N()

	if n == 0 {
		err = io.EOF
		return
	}

	actual := mh.GetShaLike()

	// if m.Verzeichnisse.Sha.IsNull() {
	m.Verzeichnisse.Sha = sha.Make(actual)
	// } else if !m.Verzeichnisse.Sha.EqualsSha(actual) &&
	// o.IncludeVerzeichnisse {
	// 	err = errors.Errorf(
	// 		"expected %q but got %q",
	// 		m.Verzeichnisse.Sha,
	// 		actual,
	// 	)
	// 	return
	// }

	c.SetMetadatei(m)

	return
}
