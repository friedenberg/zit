package zettel_printer

import (
	"io"
	"strings"
)

type stringWriterTo string

func (s stringWriterTo) WriteTo(w io.Writer) (n int64, err error) {
	var n1 int
	n1, err = io.WriteString(w, string(s))
	n = int64(n1)
	return
}

func (zp Printer) zettelBracketed(id, rev, bez io.WriterTo) string {
	sb := &strings.Builder{}

	sb.WriteString("[")

	if _, zp.Err = id.WriteTo(sb); !zp.IsEmpty() {
		zp.Wrap()
		return ""
	}

	if zp.newZettelShaSyntax {
		sb.WriteString("@")

		if _, zp.Err = rev.WriteTo(sb); !zp.IsEmpty() {
			zp.Wrap()
			return ""
		}
	}

	if zp.includeBezeichnungen {
		//TODO use bez descriptors instead
		sb.WriteString(" ")

		if _, zp.Err = bez.WriteTo(sb); !zp.IsEmpty() {
			zp.Wrap()
			return ""
		}
	}

	sb.WriteString("]")

	return sb.String()
}
