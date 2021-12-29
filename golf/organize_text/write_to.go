package organize_text

import (
	"io"
)

func (ot organizeText) WriteTo(out io.Writer) (n int64, err error) {
	w := _LineFormatNewWriter()

	w.WriteLines(_MetadateiBoundary)

	if len(ot.etiketten) > 0 {
		for _, e := range ot.etiketten {
			w.WriteFormat("* %s", e)
		}
	} else {
		w.WriteLines("*")
	}

	w.WriteLines(_MetadateiBoundary)
	w.WriteEmpty()

	n, err = w.WriteTo(out)

	var n1 int64

	if n1, err = ot.zettels.WriteTo(out); err != nil {
		err = _Error(err)
		return
	}

	n += int64(n1)

	return
}
