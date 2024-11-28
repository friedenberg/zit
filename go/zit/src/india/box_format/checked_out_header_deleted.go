package box_format

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type CheckedOutHeaderDeleted struct {
	interfaces.ConfigDryRunReader
}

func (f CheckedOutHeaderDeleted) WriteBoxHeader(
	header *string_format_writer.BoxHeader,
	co *sku.CheckedOut,
) (err error) {
	header.RightAligned = true

	if f.IsDryRun() {
		header.Value = "would delete"
	} else {
		header.Value = "deleted"
	}

	return
}
