package objekte_format

import (
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/catgut"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/ohio"
	"github.com/friedenberg/zit/src/echo/kennung"
)

type key struct {
	bytes        []byte
	includeInSha bool
}

var (
	keyAkte                         = catgut.MakeFromString("Akte")
	keyBezeichnung                  = catgut.MakeFromString("Bezeichnung")
	keyEtikett                      = catgut.MakeFromString("Etikett")
	keyGattung                      = catgut.MakeFromString("Gattung")
	keyKennung                      = catgut.MakeFromString("Kennung")
	keyKomment                      = catgut.MakeFromString("Komment")
	keyTai                          = catgut.MakeFromString("Tai")
	keyTyp                          = catgut.MakeFromString("Typ")
	keyVerzeichnisseArchiviert      = catgut.MakeFromString("Verzeichnisse-Archiviert")
	keyVerzeichnisseEtikettImplicit = catgut.MakeFromString("Verzeichnisse-Etikett-Implicit")
	keyVerzeichnisseEtikettExpanded = catgut.MakeFromString("Verzeichnisse-Etikett-Expanded")
	keyVerzeichnisseMutter          = catgut.MakeFromString("Verzeichnisse-Mutter")
	keyVerzeichnisseSha             = catgut.MakeFromString("Verzeichnisse-Sha")
)

var (
	errV4ExpectedNewline           = errors.New("expected newline")
	errV4ExpectedSpaceSeparatedKey = errors.New("expected space separated key")
	errV4EmptyKey                  = errors.New("empty key")
	errV4KeysNotSorted             = errors.New("keys not sorted")
	errV4InvalidKey                = errors.New("invalid key")
)

func (f v4) ParsePersistentMetadatei(
	r *catgut.RingBuffer,
	c ParserContext,
	o Options,
) (n int64, err error) {
	m := c.GetMetadatei()

	var (
		g gattung.Gattung
		k *kennung.Kennung2
	)

	var (
		lastKey, valBuffer       catgut.String
		line, key, val           catgut.Slice
		ok                       bool
		writeMetadateiHashString bool
		ignoreOrder              bool
	)

	mh := sha.MakeWriter(nil)
	lineNo := 0

	for {
		line, err = r.PeekUpto('\n')

		if err != nil && err != io.EOF {
			break
		}

		if line.Len() == 0 {
			break
		}

		key, val, ok = line.Cut(' ')

		if !ok {
			err = errV4ExpectedSpaceSeparatedKey
			break
		}

		if key.Len() == 0 {
			err = errV4EmptyKey
			break
		}

		if lastKey.Len() > 0 && key.Compare(lastKey.Bytes()) == 1 {
			err = errV4KeysNotSorted
			break
		}

		{
			valBuffer.Reset()
			n, err := val.WriteTo(&valBuffer)

			if n != int64(val.Len()) || err != nil {
				panic(fmt.Sprintf("failed to write val to valBuffer. N: %d, Err: %s", n, err))
			}
		}

		switch {
		case key.Equal(keyAkte.Bytes()):
			if err = m.AkteSha.SetHexBytes(valBuffer.Bytes()); err != nil {
				err = errors.Wrap(err)
				return
			}

			writeMetadateiHashString = true

		case key.Equal(keyBezeichnung.Bytes()):
			if err = m.Bezeichnung.Set(val.String()); err != nil {
				err = errors.Wrap(err)
				return
			}

			writeMetadateiHashString = true

		case key.Equal(keyEtikett.Bytes()):
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

		case key.Equal(keyGattung.Bytes()):
			if err = g.Set(val.String()); err != nil {
				err = errors.Wrap(err)
				return
			}

		case key.Equal(keyKennung.Bytes()):
			k = kennung.GetKennungPool().Get()
			defer kennung.GetKennungPool().Put(k)

			if err = k.SetWithGattung(val.String(), g); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = c.SetKennungLike(k); err != nil {
				err = errors.Wrap(err)
				return
			}

		case key.Equal(keyTai.Bytes()):
			if err = m.Tai.Set(val.String()); err != nil {
				err = errors.Wrap(err)
				return
			}

			writeMetadateiHashString = true

		case key.Equal(keyTyp.Bytes()):
			if err = m.Typ.Set(val.String()); err != nil {
				err = errors.Wrap(err)
				return
			}

			writeMetadateiHashString = true

		case key.Equal(keyVerzeichnisseArchiviert.Bytes()):
			if err = m.Verzeichnisse.Archiviert.Set(val.String()); err != nil {
				err = errors.Wrap(err)
				return
			}

		case key.Equal(keyVerzeichnisseEtikettImplicit.Bytes()):
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

			if err = m.Verzeichnisse.AddEtikettImplicitPtr(e); err != nil {
				err = errors.Wrap(err)
				return
			}

		case key.Equal(keyVerzeichnisseEtikettExpanded.Bytes()):
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

			if err = m.Verzeichnisse.AddEtikettExpandedPtr(e); err != nil {
				err = errors.Wrap(err)
				return
			}

		case key.Equal(keyVerzeichnisseMutter.Bytes()):
			if err = m.Verzeichnisse.Mutter.SetHexBytes(val.Bytes()); err != nil {
				err = errors.Wrap(err)
				return
			}

			writeMetadateiHashString = true

		case key.Equal(keyVerzeichnisseSha.Bytes()):
			if err = m.Verzeichnisse.Sha.SetHexBytes(val.Bytes()); err != nil {
				err = errors.Wrap(err)
				return
			}

		case key.Equal(keyKomment.Bytes()):
			ignoreOrder = true
			m.Comments = append(m.Comments, val.String())

		default:
			err = errV4InvalidKey
		}

		// Key Space Value Newline
		thisN := int64(key.Len() + 1 + val.Len() + 1)
		n += thisN

		if !ignoreOrder {
			lastKey.Reset()
			n, err = key.WriteTo(&lastKey)

			if n != int64(key.Len()) || err != nil {
				panic(
					fmt.Sprintf(
						"failed to write everything to lastKey. N: %d, err: %s",
						n,
						err,
					),
				)
			}
		}

		lineNo++

		r.AdvanceRead(int(thisN))

		if writeMetadateiHashString {
			if _, err = ohio.WriteKeySpaceValueNewlineWritersTo(
				mh,
				key,
				val,
			); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		writeMetadateiHashString = false
		ignoreOrder = false
	}

	if n == 0 {
		if err == nil {
			err = io.EOF
		}

		return
	}

	actual := mh.GetShaLike()

	// if m.Verzeichnisse.Sha.IsNull() {
	if err = m.Verzeichnisse.Sha.SetShaLike(actual); err != nil {
		err = errors.Wrap(err)
		return
	}

	// } else if !m.Verzeichnisse.Sha.EqualsSha(actual) &&
	// o.IncludeVerzeichnisse {
	// 	err = errors.Errorf(
	// 		"expected %q but got %q",
	// 		m.Verzeichnisse.Sha,
	// 		actual,
	// 	)
	// 	return
	// }

	return
}
