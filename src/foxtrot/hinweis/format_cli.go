package hinweis

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/format"
	"github.com/friedenberg/zit/src/schnittstellen"
)

// kopf/schwanz
func MakeCliFormat(
	cw format.FuncColorWriter,
	a schnittstellen.FuncAbbreviateKorper,
	maxKopf int,
	maxSchwanz int,
) format.FormatWriterFunc[Hinweis] {
	return func(w io.Writer, h Hinweis) (n int64, err error) {
		if a != nil {
			var v string

			if v, err = a(h); err != nil {
				err = errors.Wrap(err)
				return
			}

			if v != "" {
				if err = h.Set(v); err != nil {
					err = errors.Wrap(err)
					return
				}
			} else {
				errors.Todo("empty hinweis abbr")
			}
		}

		p1, p2 := h.AlignedParts(maxKopf, maxSchwanz)

		return format.Write(
			w,
			cw(format.MakeFormatString("%s", p1), format.ColorTypePointer),
			format.MakeFormatString("/"),
			cw(format.MakeFormatString("%s", p2), format.ColorTypePointer),
		)
	}
}
