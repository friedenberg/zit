package sku_fmt

import (
	"fmt"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func String(o *sku.Transacted) (str string) {
	str = fmt.Sprintf(
		"%s %s %s %s %s",
		o.GetTai(),
		o.GetGenre(),
		o.GetObjectId(),
		o.GetObjectSha(),
		o.GetBlobSha(),
	)

	return
}

func StringMetadatei(o *sku.Transacted) (str string) {
	sb := &strings.Builder{}

	sb.WriteString(o.GetTai().String())

	sb.WriteString(" ")
	sb.WriteString(o.GetGenre().GetGenreString())

	sb.WriteString(" ")
	sb.WriteString(o.GetObjectId().String())

	sb.WriteString(" ")
	sb.WriteString(o.GetBlobSha().String())

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
			iter.StringDelimiterSeparated[ids.Tag](
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

	return sb.String()
}
