package object_metadata

import (
	"fmt"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

func StringSansTai(o *Metadata) (str string) {
	sb := &strings.Builder{}

	sb.WriteString(" ")
	sb.WriteString(o.Blob.String())

	m := o.GetMetadata()

	t := m.GetType()

	if !t.IsEmpty() {
		sb.WriteString(" ")
		sb.WriteString(ids.FormattedString(m.GetType()))
	}

	es := m.GetTags()

	if es.Len() > 0 {
		sb.WriteString(" ")
		sb.WriteString(
			quiter.StringDelimiterSeparated[ids.Tag](
				" ",
				m.GetTags(),
			),
		)
	}

	b := m.Description

	if !b.IsEmpty() {
		sb.WriteString(" ")
		sb.WriteString("\"" + b.String() + "\"")
	}

	for _, field := range m.Fields {
		sb.WriteString(" ")
		fmt.Fprintf(sb, "%q=%q", field.Key, field.Value)
	}

	return sb.String()
}
