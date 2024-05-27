package objekte_format

import (
	"io"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/src/delta/sha"
	"code.linenisgreat.com/zit/src/echo/kennung"
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

	if !m.Akte.IsNull() {
		n1, err = ohio.WriteKeySpaceValueNewlineString(
			mw,
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
			mw,
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

	for _, e := range iter.SortedValues[kennung.Etikett](es) {
		n1, err = ohio.WriteKeySpaceValueNewlineString(
			mw,
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
		c.GetKennung().GetGattung().GetGattungString(),
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
			mw,
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
			mw,
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
			mw,
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

			for _, e := range iter.SortedValues[kennung.Etikett](
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

			for _, e := range iter.SortedValues[kennung.Etikett](
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

	if !m.Mutter().IsNull() && !o.ExcludeMutter {
		n1, err = ohio.WriteKeySpaceValueNewlineString(
			mw,
			keyShasMutterMetadateiKennungMutter.String(),
			m.Mutter().String(),
		)
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if o.PrintFinalSha {
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
			keySha.String(),
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
