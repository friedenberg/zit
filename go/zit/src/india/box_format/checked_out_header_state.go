package box_format

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type CheckedOutHeaderState struct{}

func (f CheckedOutHeaderState) WriteBoxHeader(
	header *string_format_writer.BoxHeader,
	co *sku.CheckedOut,
) (err error) {
	state := co.GetState()
	stateString := state.String()
	header.RightAligned = true
	header.Value = stateString

	return
}
