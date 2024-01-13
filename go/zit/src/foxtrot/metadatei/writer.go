package metadatei

import (
	"bufio"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type Writer struct {
	Metadatei, Akte io.WriterTo
}

func (w1 Writer) WriteTo(w2 io.Writer) (n int64, err error) {
	w := bufio.NewWriter(w2)
	defer errors.Deferred(&err, w.Flush)

	var n1 int64
	var n2 int

	hasMetadateiContent := true

	if mwt, ok := w1.Metadatei.(MetadateiWriterTo); ok {
		hasMetadateiContent = mwt.HasMetadateiContent()
	}

	if w1.Metadatei != nil && hasMetadateiContent {
		n2, err = w.WriteString(Boundary + "\n")
		n += int64(n2)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n1, err = w1.Metadatei.WriteTo(w)
		n += n1

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		w.WriteString(Boundary + "\n")
		n += n1

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		if w1.Akte != nil {
			w.WriteString("\n")
			n += n1

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	if w1.Akte != nil {
		n1, err = w1.Akte.WriteTo(w)
		n += n1

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
