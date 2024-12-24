package sku_fmt_debug

import (
	"fmt"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func StringTaiGenreObjectIdShaBlob(o *sku.Transacted) (str string) {
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

func StringObjectIdBlobMetadataSansTai(o *sku.Transacted) (str string) {
	str = fmt.Sprintf(
		"%s %s %s",
		o.GetObjectId(),
		o.GetBlobSha(),
		StringMetadataSansTai(o),
	)

	return
}

func StringMetadataTai(o *sku.Transacted) (str string) {
	return fmt.Sprintf("%s %s", o.GetTai(), StringMetadataSansTai(o))
}

func StringMetadataSansTai(o *sku.Transacted) (str string) {
	sb := &strings.Builder{}

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
