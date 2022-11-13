package hinweis

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
)

// kopf/schwanz
func MakeCliFormat(
	a Abbr,
	maxKopf int,
	maxSchwanz int,
) collections.WriterFuncFormat[Hinweis] {
	return func(w io.Writer, h *Hinweis) (n int64, err error) {
		if a != nil {
			if *h, err = a.AbbreviateHinweis(*h); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		var n1 int

		if n1, err = io.WriteString(w, h.Aligned(maxKopf, maxSchwanz)); err != nil {
			err = errors.Wrap(err)
			return
		}

		n += int64(n1)

		return
	}
}
