package hinweis

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/format"
)

// kopf/schwanz
func MakeCliFormat(
	cw format.FuncColorWriter,
	a Abbr,
	maxKopf int,
	maxSchwanz int,
) format.FormatWriterFunc[Hinweis] {
	return func(w io.Writer, h *Hinweis) (n int64, err error) {
		if a != nil {
			if *h, err = a.AbbreviateHinweis(*h); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		return format.Write(
			w,
			cw(format.MakeFormatString(h.Aligned(maxKopf, maxSchwanz)), format.ColorTypeConstant),
		)
	}
}
