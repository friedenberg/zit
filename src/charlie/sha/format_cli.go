package sha

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/format"
)

// sha
func MakeCliFormat(
	cw format.FuncColorWriter,
	a Abbr,
) format.FormatWriterFunc[Sha] {
	return func(w io.Writer, s *Sha) (n int64, err error) {
		v := s.String()

		if a != nil {
			if v, err = a.AbbreviateSha(*s); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		return format.Write(
			w,
			cw(format.MakeFormatString(v), format.ColorTypeConstant),
		)
	}
}
