package box_format

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type CheckedOutHeaderUserTai struct{}

func (f CheckedOutHeaderUserTai) WriteBoxHeader(
	header *string_format_writer.BoxHeader,
	co *sku.CheckedOut,
) (err error) {
	external := co.GetSkuExternal()
	t := external.GetTai()
	header.Value = t.Format(string_format_writer.StringFormatDateTime)

	return
}
