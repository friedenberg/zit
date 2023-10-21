package objekte_format

import (
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
	var n1 int

	if !m.AkteSha.IsNull() {
		n1, err = ohio.WriteKeySpaceValueNewline(
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

		n1, err = ohio.WriteKeySpaceValueNewline(
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
		n1, err = ohio.WriteKeySpaceValueNewline(
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

	n1, err = ohio.WriteKeySpaceValueNewline(
		w,
		"Gattung",
		c.GetKennungLike().GetGattung().GetGattungString(),
	)
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = ohio.WriteKeySpaceValueNewline(
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
		n1, err = ohio.WriteKeySpaceValueNewline(
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
		n1, err = ohio.WriteKeySpaceValueNewline(
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
			n1, err = ohio.WriteKeySpaceValueNewline(
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

		if m.Verzeichnisse.ExpandedEtiketten != nil {
			k := fmt.Sprintf(
				"Verzeichnisse-%s-Expanded",
				gattung.Etikett.String(),
			)
			for _, e := range iter.SortedValues[kennung.Etikett](m.Verzeichnisse.ExpandedEtiketten) {
				n1, err = ohio.WriteKeySpaceValueNewline(
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

		if m.Verzeichnisse.ImplicitEtiketten != nil {
			k := fmt.Sprintf(
				"Verzeichnisse-%s-Implicit",
				gattung.Etikett.String(),
			)

			for _, e := range iter.SortedValues[kennung.Etikett](m.Verzeichnisse.ImplicitEtiketten) {
				n1, err = ohio.WriteKeySpaceValueNewline(
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

		if !m.Verzeichnisse.Mutter.IsNull() {
			n1, err = ohio.WriteKeySpaceValueNewline(
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

		n1, err = ohio.WriteKeySpaceValueNewline(
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

	etiketten := kennung.MakeEtikettMutableSet()
	var etikettenExpanded, etikettenImplicit kennung.EtikettMutableSet

	if o.IncludeVerzeichnisse {
		etikettenExpanded = kennung.MakeEtikettMutableSet()
		etikettenImplicit = kennung.MakeEtikettMutableSet()
	}

	var (
		g gattung.Gattung
		k kennung.Kennung2
	)

	dr := ohio.MakeDelimReader('\n', r1)
	defer ohio.PutDelimReader(dr)

	var (
		lastKey string
		key     string
		val     string
	)

	mh := sha.MakeWriter(nil)

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

		writeMetadateiHashString := false

		switch key {
		case "Akte":
			if err = m.AkteSha.Set(val); err != nil {
				err = errors.Wrap(err)
				return
			}

			writeMetadateiHashString = true

		case "Bezeichnung":
			if err = m.Bezeichnung.Set(val); err != nil {
				err = errors.Wrap(err)
				return
			}

			writeMetadateiHashString = true

		case "Etikett":
			if err = iter.AddString[kennung.Etikett, *kennung.Etikett](
				etiketten,
				val,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			writeMetadateiHashString = true

		case "Gattung":
			if err = g.Set(val); err != nil {
				err = errors.Wrap(err)
				return
			}

		case "Kennung":
			if err = k.SetWithGattung(val, g); err != nil {
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

			writeMetadateiHashString = true

		case "Typ":
			if err = m.Typ.Set(val); err != nil {
				err = errors.Wrap(err)
				return
			}

			writeMetadateiHashString = true

		case "Verzeichnisse-Archiviert":
			if err = m.Verzeichnisse.Archiviert.Set(val); err != nil {
				err = errors.Wrap(err)
				return
			}

		case "Verzeichnisse-Etikett-Implicit":
			if !o.IncludeVerzeichnisse {
				err = errors.Errorf(
					"format specifies not to include Verzeichnisse but found %q",
					key,
				)
				return
			}

			if err = iter.AddString[kennung.Etikett, *kennung.Etikett](
				etikettenImplicit,
				val,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

		case "Verzeichnisse-Etikett-Expanded":
			if !o.IncludeVerzeichnisse {
				err = errors.Errorf(
					"format specifies not to include Verzeichnisse but found %q",
					key,
				)
				return
			}

			if err = iter.AddString[kennung.Etikett, *kennung.Etikett](
				etikettenExpanded,
				val,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

		case "Verzeichnisse-Mutter":
			if err = m.Verzeichnisse.Mutter.Set(val); err != nil {
				err = errors.Wrap(err)
				return
			}

			writeMetadateiHashString = true

		case "Verzeichnisse-Sha":
			if err = m.Verzeichnisse.Sha.Set(val); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		lastKey = key

		if !writeMetadateiHashString {
			continue
		}

		if _, err = ohio.WriteKeySpaceValueNewline(mh, key, val); err != nil {
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

	m.Etiketten = etiketten

	if o.IncludeVerzeichnisse {
		m.Verzeichnisse.ImplicitEtiketten = etikettenImplicit
		m.Verzeichnisse.ExpandedEtiketten = etikettenExpanded
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
