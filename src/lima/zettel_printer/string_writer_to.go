package zettel_printer

import "io"

type stringWriterTo string

func (s stringWriterTo) WriteTo(w io.Writer) (n int64, err error) {
	var n1 int
	n1, err = io.WriteString(w, string(s))
	n = int64(n1)
	return
}
