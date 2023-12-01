package objekte_format

import (
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/catgut"
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

var (
	errV4ExpectedNewline           = errors.New("expected newline")
	errV4ExpectedSpaceSeparatedKey = errors.New("expected space separated key")
	errV4EmptyKey                  = errors.New("empty key")
	errV4KeysNotSorted             = errors.New("keys not sorted")
	errV4InvalidKey                = errors.New("invalid key")
)

func (f v4) ParsePersistentMetadatei(
	r *ohio_ring_buffer2.RingBuffer,
	c ParserContext,
	o Options,
) (n int64, err error) {
	m := c.GetMetadatei()

	var (
		g gattung.Gattung
		k *kennung.Kennung2
	)

	var (
		lastKey, valBuffer catgut.String
		line, key, val     ohio_ring_buffer2.Slice
		ok                 bool
	)

	mh := sha.MakeWriter(nil)
	lineNo := 0

	for {
		line, ok, err = r.PeekUpto('\n')

		if err != nil && err != io.EOF {
			break
		}

		if !ok && err != io.EOF {
			err = errV4ExpectedNewline
			return
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

		writeMetadateiHashString := false

		if key.Equal(keyAkte) {
			if err = m.AkteSha.SetHexBytes(valBuffer.Bytes()); err != nil {
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

			if err = m.Verzeichnisse.AddEtikettImplicitPtr(e); err != nil {
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

			if err = m.Verzeichnisse.AddEtikettExpandedPtr(e); err != nil {
				err = errors.Wrap(err)
				return
			}
		} else if key.Equal(keyVerzeichnisseMutter) {
			if err = m.Verzeichnisse.Mutter.SetHexBytes(val.Bytes()); err != nil {
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
			err = errV4InvalidKey
			break
		}

		// Key Space Value Newline
		thisN := int64(key.Len() + 1 + val.Len() + 1)
		n += thisN

		{
			lastKey.Reset()
			n, err := key.WriteTo(&lastKey)

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

	return
}
