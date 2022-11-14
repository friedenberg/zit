package hinweis

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/format"
)

// kopf/schwanz
func MakeCliFormat(
	cw format.FuncColorWriter,
	a Abbr,
	maxKopf int,
	maxSchwanz int,
) format.FormatWriterFunc[Hinweis] {
	return func(w io.Writer, h *Hinweis) (n int64, err error) {
		h1 := *h

		if a != nil {
			if h1, err = a.AbbreviateHinweis(h1); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		return format.Write(
			w,
      //TODO do not use color for slash
			cw(format.MakeFormatString(h1.Aligned(maxKopf, maxSchwanz)), format.ColorTypePointer),
		)
	}
}
