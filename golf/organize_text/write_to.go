package organize_text

import (
	"io"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/line_format"
	"github.com/friedenberg/zit/foxtrot/zettel_formats"
)

func (ot organizeText) WriteTo(out io.Writer) (n int64, err error) {
	w := line_format.NewWriter()

	w.WriteLines(zettel_formats.MetadateiBoundary)

	if len(ot.etiketten) > 0 {
		for _, e := range ot.etiketten {
			w.WriteFormat("* %s", e)
		}
	} else {
		w.WriteLines("*")
	}

	w.WriteLines(zettel_formats.MetadateiBoundary)
	w.WriteEmpty()

	n, err = w.WriteTo(out)

	var n1 int64

	if n1, err = ot.zettels.WriteTo(out); err != nil {
		err = errors.Error(err)
		return
	}

	n += int64(n1)

	return
}
