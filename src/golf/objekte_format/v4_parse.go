package objekte_format

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/ohio_ring_buffer2"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/ohio"
	"github.com/friedenberg/zit/src/echo/kennung"
)

type key struct {
	bytes        []byte
	includeInSha bool
}

var (
	keyAkte                         = []byte("Akte")
	keyBezeichnung                  = []byte("Bezeichnung")
	keyEtikett                      = []byte("Etikett")
	keyGattung                      = []byte("Gattung")
	keyKennung                      = []byte("Kennung")
	keyTai                          = []byte("Tai")
	keyTyp                          = []byte("Typ")
	keyVerzeichnisseArchiviert      = []byte("Verzeichnisse-Archiviert")
	keyVerzeichnisseEtikettImplicit = []byte("Verzeichnisse-Etikett-Implicit")
	keyVerzeichnisseEtikettExpanded = []byte("Verzeichnisse-Etikett-Expanded")
	keyVerzeichnisseMutter          = []byte("Verzeichnisse-Mutter")
	keyVerzeichnisseSha             = []byte("Verzeichnisse-Sha")
)

func (f v4) ParsePersistentMetadatei(
	r *ohio_ring_buffer2.RingBuffer,
	c ParserContext,
	o Options,
) (n int64, err error) {
	m := c.GetMetadatei()

	var (
		g gattung.Gattung
		k kennung.Kennung2
	)

	var (
		lastKey        []byte
		line, key, val ohio_ring_buffer2.Slice
		ok             bool
	)

	mh := sha.MakeWriter(nil)
	lineNo := 0

	for {
		line, ok, err = r.PeekUpto('\n')

		if err != nil && err != io.EOF {
			break
		}

		if !ok && err != io.EOF {
			err = errors.Errorf("expected a newline terminated string %q", line)
			return
		}

		if line.Len() == 0 {
			break
		}

		key, val, ok = line.Cut(' ')

		if !ok {
			err = errors.Errorf("expected space-separated key-value but got %q", line)
			break
		}

		if key.Len() == 0 {
			err = errors.Errorf("empty key at line %d", lineNo)
			break
		}

		if len(lastKey) > 0 && key.Compare(lastKey) == 1 {
			err = errors.Errorf("keys not sorted: last: %q, current: %q", lastKey, key)
			break
		}

		writeMetadateiHashString := false

		if key.Equal(keyAkte) {
			if err = m.AkteSha.Set(val.String()); err != nil {
				err = errors.Wrap(err)
				return
			}

			writeMetadateiHashString = true

		} else if key.Equal(keyBezeichnung) {
			if err = m.Bezeichnung.Set(val.String()); err != nil {
				err = errors.Wrap(err)
				return
			}

			writeMetadateiHashString = true

		} else if key.Equal(keyEtikett) {
			e := kennung.GetEtikettPool().Get()

			if err = e.Set(val.String()); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = m.AddEtikettPtr(e); err != nil {
				err = errors.Wrap(err)
				return
			}

			writeMetadateiHashString = true

		} else if key.Equal(keyGattung) {
			if err = g.Set(val.String()); err != nil {
				err = errors.Wrap(err)
				return
			}
		} else if key.Equal(keyKennung) {
			if err = k.SetWithGattung(val.String(), g); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = c.SetKennungLike(k); err != nil {
				err = errors.Wrap(err)
				return
			}

		} else if key.Equal(keyTai) {
			if err = m.Tai.Set(val.String()); err != nil {
				err = errors.Wrap(err)
				return
			}

			writeMetadateiHashString = true

		} else if key.Equal(keyTyp) {
			if err = m.Typ.Set(val.String()); err != nil {
				err = errors.Wrap(err)
				return
			}

			writeMetadateiHashString = true

		} else if key.Equal(keyVerzeichnisseArchiviert) {
			if err = m.Verzeichnisse.Archiviert.Set(val.String()); err != nil {
				err = errors.Wrap(err)
				return
			}
		} else if key.Equal(keyVerzeichnisseEtikettImplicit) {
			if !o.IncludeVerzeichnisse {
				err = errors.Errorf(
					"format specifies not to include Verzeichnisse but found %q",
					key,
				)
				return
			}

			e := kennung.GetEtikettPool().Get()

			if err = e.Set(val.String()); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = m.Verzeichnisse.GetImplicitEtikettenMutable().AddPtr(e); err != nil {
				err = errors.Wrap(err)
				return
			}

		} else if key.Equal(keyVerzeichnisseEtikettExpanded) {
			if !o.IncludeVerzeichnisse {
				err = errors.Errorf(
					"format specifies not to include Verzeichnisse but found %q",
					key,
				)
				return
			}

			e := kennung.GetEtikettPool().Get()

			if err = e.Set(val.String()); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = m.Verzeichnisse.GetExpandedEtikettenMutable().AddPtr(e); err != nil {
				err = errors.Wrap(err)
				return
			}
		} else if key.Equal(keyVerzeichnisseMutter) {
			if err = m.Verzeichnisse.Mutter.Set(val.String()); err != nil {
				err = errors.Wrap(err)
				return
			}

			writeMetadateiHashString = true

		} else if key.Equal(keyVerzeichnisseSha) {
			if err = m.Verzeichnisse.Sha.Set(val.String()); err != nil {
				err = errors.Wrap(err)
				return
			}
		} else {
			err = errors.Errorf("not a valid key: %q", key)
			break
		}

		// Key Space Value Newline
		thisN := int64(key.Len() + 1 + val.Len() + 1)
		n += thisN

		lastKey = []byte(key.String())

		lineNo++

		r.AdvanceRead(int(thisN))

		if !writeMetadateiHashString {
			continue
		}

		if _, err = ohio.WriteKeySpaceValueNewlineWritersTo(
			mh,
			key,
			val,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		writeMetadateiHashString = false
	}

	if n == 0 {
		if err == nil {
			err = io.EOF
		}

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
