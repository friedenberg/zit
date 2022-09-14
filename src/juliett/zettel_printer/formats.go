package zettel_printer

import (
	"strings"
)

func (zp Printer) zettelBracketed(id, rev, bez string) string {
	sb := &strings.Builder{}

	sb.WriteString("[")
	sb.WriteString(id)

	if zp.newZettelShaSyntax {
		sb.WriteString("@")
		sb.WriteString(rev)
	}

	if zp.includeBezeichnungen {
		//TODO use bez descriptors instead
		sb.WriteString(" ")
		sb.WriteString(bez)
	}

	sb.WriteString("]")

	return sb.String()
}
