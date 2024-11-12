package box_format

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

// TODO modify to include diff state depending on checked out state
type CheckedOutHeaderState struct{}

func (f CheckedOutHeaderState) WriteBoxHeader(
	header *string_format_writer.BoxHeader,
	co *sku.CheckedOut,
) (err error) {
	header.RightAligned = true

	state := co.GetState()
	stateString := state.String()

	switch state {
	// case checked_out_state.JustCheckedOut:
	// 	header.Value = stateString

	case checked_out_state.CheckedOut:
		if co.GetSku().Metadata.EqualsSansTai(&co.GetSkuExternal().GetSku().Metadata) {
			header.Value = string_format_writer.StringSame
		} else {
			header.Value = string_format_writer.StringChanged
		}

		// 	case checked_out_state.Untracked:
		// 		header.Value = stateString
		// 	case checked_out_state.Recognized:
		// 		header.Value = stateString

		// 	case checked_out_state.Conflicted:
		// 		header.Value = stateString

	default:
		header.Value = stateString
	}

	return
}
