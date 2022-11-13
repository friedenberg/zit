package bezeichnung

import (
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
)

func MakeCliFormat() collections.WriterFuncFormat[Bezeichnung] {
	return func(w io.Writer, b1 *Bezeichnung) (n int64, err error) {
		b := b1.value

		switch {
		case len(b) > 66:
			b = b[:66] + "…"
		}

		var n1 int

		if n1, err = io.WriteString(w, fmt.Sprintf("\"%s\"", b)); err != nil {
			n = int64(n1)
			err = errors.Wrap(err)
			return
		}

		n = int64(n1)

		return
	}
}
