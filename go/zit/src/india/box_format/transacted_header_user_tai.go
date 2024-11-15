package box_format

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type TransactedHeaderUserTai struct{}

func (f TransactedHeaderUserTai) WriteBoxHeader(
	header *string_format_writer.BoxHeader,
	sk *sku.Transacted,
) (err error) {
	t := sk.GetTai()
	header.RightAligned = true
	header.Value = t.Format(string_format_writer.StringFormatDateTime)

	return
}
