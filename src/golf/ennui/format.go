package ennui

import (
	"fmt"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/catgut"
	"github.com/friedenberg/zit/src/delta/ohio"
	"github.com/friedenberg/zit/src/echo/kennung"
)

type format struct {
	key  string
	keys []*catgut.String
}

func (f format) printKeys(
	w io.Writer,
	m *Metadatei,
) (n int64, err error) {
	var n1 int64

	for _, k := range f.keys {
		n1, err = f.printKey(w, m, k)
		n += n1

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (f format) printKey(
	w io.Writer,
	m *Metadatei,
	key *catgut.String,
) (n int64, err error) {
	var (
		n1 int
		n2 int64
	)

	switch key {
	case keyAkte:
		if !m.Akte.IsNull() {
			n2, err = catgut.WriteKeySpaceValueNewline(
				w,
				keyAkte,
				&m.Akte,
			)
			n += int64(n2)

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	case keyBezeichnung:
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

	case keyEtikett:
		es := m.GetEtiketten()

		for _, e := range iter.SortedValues[kennung.Etikett](es) {
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

	case keyMutter:

		if !m.Mutter.IsNull() {
			n1, err = ohio.WriteKeySpaceValueNewlineString(
				w,
				keyMutter.String(),
				m.Mutter.String(),
			)
			n += int64(n1)

			if err != nil {
				err = errors.Wrap(err)
				return
			}

			n += int64(n1)

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}

	case keyTai:
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

	case keyTyp:
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

	default:
		panic(fmt.Sprintf("unsupported key: %s", key))
	}

	return
}
