package standort

import (
	"io"
	"path/filepath"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/format"
)

func (s Standort) MakeWriterRelativePath(
	p string,
) format.WriterFunc {
	return func(w io.Writer) (n int64, err error) {
		var n1 int

		{
			// if p, err = filepath.Rel(s.cwd, p); err != nil {
			// 	err = errors.Wrap(err)
			// 	return
			// }

			p1, _ := filepath.Rel(s.cwd, p)

			if p1 != "" {
				p = p1
			}
		}

		if n1, err = io.WriteString(w, p); err != nil {
			n = int64(n1)
			err = errors.Wrap(err)
			return
		}

		n = int64(n1)

		return
	}
}
