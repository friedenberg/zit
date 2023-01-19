package sha_cli_format

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/format"
)

// sha
func MakeCliFormat(
	cw format.FuncColorWriter,
	a schnittstellen.FuncAbbreviateValue,
) schnittstellen.FuncWriterFormat[schnittstellen.Sha] {
	return func(w io.Writer, s schnittstellen.Sha) (n int64, err error) {
		v := s.String()

		if a != nil {
			var v1 string

			if v1, err = a(s); err != nil {
				err = errors.Wrap(err)
				return
			}

			if v1 != "" {
				v = v1
			} else {
				errors.Todo("abbreviate sha produced empty string")
			}
		}

		return format.Write(
			w,
			cw(format.MakeFormatString(v), format.ColorTypeConstant),
		)
	}
}
