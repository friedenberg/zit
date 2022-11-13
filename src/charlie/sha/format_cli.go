package sha

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
)

// sha
func MakeCliFormat(
	a Abbr,
) collections.WriterFuncFormat[Sha] {
	return func(w io.Writer, s *Sha) (n int64, err error) {
		v := s.String()

		if a != nil {
			if v, err = a.AbbreviateSha(*s); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		var n1 int

		if n1, err = io.WriteString(w, v); err != nil {
			err = errors.Wrap(err)
			return
		}

		n += int64(n1)

		return
	}
}
