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

		p1, p2 := h1.AlignedParts(maxKopf, maxSchwanz)

		return format.Write(
			w,
			cw(format.MakeFormatString("%s", p1), format.ColorTypePointer),
			format.MakeFormatString("/"),
			cw(format.MakeFormatString("%s", p2), format.ColorTypePointer),
		)
	}
}
