package object_inventory_format

import (
	"fmt"
	"io"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type v5 struct{}

func (f v5) FormatPersistentMetadatei(
	w1 io.Writer,
	c FormatterContext,
	o Options,
) (n int64, err error) {
	w := ohio.GetPoolBufioWriter().Get()
	defer ohio.GetPoolBufioWriter().Put(w)

	w.Reset(w1)
	defer errors.DeferredFlusher(&err, w)

	m := c.GetMetadatei()

	var (
		n1 int
		n2 int64
	)

	if !m.Akte.IsNull() {
		n1, err = ohio.WriteKeySpaceValueNewlineString(
			w,
			keyAkte.String(),
			m.Akte.String(),
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
			w,
			keyBezeichnung.String(),
			line,
		)
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	es := m.GetEtiketten()

	for _, e := range iter.SortedValues[ids.Tag](es) {
		if e.IsVirtual() {
			continue
		}

		n1, err = ohio.WriteKeySpaceValueNewlineString(
			w,
			keyEtikett.String(),
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
		keyGattung.String(),
		c.GetKennung().GetGenre().GetGenreString(),
	)
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = ohio.WriteKeySpaceValueNewlineString(
		w,
		keyKennung.String(),
		c.GetKennung().String(),
	)
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, k := range m.Comments {
		n1, err = ohio.WriteKeySpaceValueNewlineString(
			w,
			keyKomment.String(),
			k,
		)
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if o.Tai {
		n1, err = ohio.WriteKeySpaceValueNewlineString(
			w,
			keyTai.String(),
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
			w,
			keyTyp.String(),
			m.GetTyp().String(),
		)
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if o.Verzeichnisse {
		if m.Verzeichnisse.Schlummernd.Bool() {
			n1, err = ohio.WriteKeySpaceValueNewlineString(
				w,
				keyVerzeichnisseArchiviert.String(),
				m.Verzeichnisse.Schlummernd.String(),
			)
			n += int64(n1)

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		if m.Verzeichnisse.GetExpandedEtiketten().Len() > 0 {
			k := keyVerzeichnisseEtikettExpanded.String()

			for _, e := range iter.SortedValues[ids.Tag](
				m.Verzeichnisse.GetExpandedEtiketten(),
			) {
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
			k := keyVerzeichnisseEtikettImplicit.String()

			for _, e := range iter.SortedValues[ids.Tag](
				m.Verzeichnisse.GetImplicitEtiketten(),
			) {
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
	}

	n1, err = writeShaKeyIfNotNull(
		w,
		keyShasMutterMetadateiKennungMutter,
		m.Mutter(),
	)

	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

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
		keySha.String(),
		m.Sha().String(),
	)

	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f v5) ParsePersistentMetadatei(
	r *catgut.RingBuffer,
	c ParserContext,
	o Options,
) (n int64, err error) {
	m := c.GetMetadatei()

	var (
		g genres.Genre
		k *ids.ObjectId
	)

	var (
		valBuffer      catgut.String
		line, key, val catgut.Slice
		ok             bool
	)

	lineNo := 0

	for {
		line, err = r.PeekUpto('\n')

		if errors.IsNotNilAndNotEOF(err) {
			break
		}

		if line.Len() == 0 {
			break
		}

		key, val, ok = line.Cut(' ')

		if !ok {
			err = makeErrWithBytes(ErrV4ExpectedSpaceSeparatedKey, line.Bytes())
			break
		}

		if key.Len() == 0 {
			err = makeErrWithBytes(errV4EmptyKey, line.Bytes())
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
			if err = m.Akte.SetHexBytes(valBuffer.Bytes()); err != nil {
				err = errors.Wrap(err)
				return
			}

		case key.Equal(keyBezeichnung.Bytes()):
			if err = m.Bezeichnung.Set(val.String()); err != nil {
				err = errors.Wrap(err)
				return
			}

		case key.Equal(keyEtikett.Bytes()):
			e := ids.GetTagPool().Get()

			if err = e.Set(val.String()); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = m.AddEtikettPtr(e); err != nil {
				err = errors.Wrap(err)
				return
			}

		case key.Equal(keyGattung.Bytes()):
			if err = g.Set(val.String()); err != nil {
				err = errors.Wrap(err)
				return
			}

		case key.Equal(keyKennung.Bytes()):
			k = ids.GetObjectIdPool().Get()
			defer ids.GetObjectIdPool().Put(k)

			if err = k.SetWithGenre(val.String(), g); err != nil {
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

		case key.Equal(keyTyp.Bytes()):
			if err = m.Typ.Set(val.String()); err != nil {
				err = errors.Wrap(err)
				return
			}

		case key.Equal(keyVerzeichnisseArchiviert.Bytes()):
			if err = m.Verzeichnisse.Schlummernd.Set(val.String()); err != nil {
				err = errors.Wrap(err)
				return
			}

		case key.Equal(keyVerzeichnisseEtikettImplicit.Bytes()):
			if !o.Verzeichnisse {
				err = errors.Errorf(
					"format specifies not to include Verzeichnisse but found %q",
					key,
				)
				return
			}

			e := ids.GetTagPool().Get()

			if err = e.Set(val.String()); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = m.Verzeichnisse.AddEtikettImplicitPtr(e); err != nil {
				err = errors.Wrap(err)
				return
			}

		case key.Equal(keyVerzeichnisseEtikettExpanded.Bytes()):
			if !o.Verzeichnisse {
				err = errors.Errorf(
					"format specifies not to include Verzeichnisse but found %q",
					key,
				)
				return
			}

			e := ids.GetTagPool().Get()

			if err = e.Set(val.String()); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = m.Verzeichnisse.AddEtikettExpandedPtr(e); err != nil {
				err = errors.Wrap(err)
				return
			}

		case key.Equal(keyShasMutterMetadateiKennungMutter.Bytes()):
			if err = m.Mutter().SetHexBytes(val.Bytes()); err != nil {
				err = errors.Wrap(err)
				return
			}

		case key.Equal(keySha.Bytes()):
			if err = m.Sha().SetHexBytes(val.Bytes()); err != nil {
				err = errors.Wrap(err)
				return
			}

		case key.Equal(keyKomment.Bytes()):
			m.Comments = append(m.Comments, val.String())

		default:
			err = errV4InvalidKey
		}

		// Key Space Value Newline
		thisN := int64(key.Len() + 1 + val.Len() + 1)
		n += thisN

		lineNo++

		r.AdvanceRead(int(thisN))
	}

	if n == 0 {
		if err == nil {
			err = io.EOF
		}

		return
	}

	return
}
