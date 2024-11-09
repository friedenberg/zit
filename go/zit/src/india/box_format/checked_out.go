package box_format

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/options_print"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func MakeBoxCheckedOut(
	co string_format_writer.ColorOptions,
	options options_print.V0,
	fieldsFormatWriter interfaces.StringFormatWriter[string_format_writer.Box],
	abbr ids.Abbr,
	fsItemReadWriter sku.FSItemReadWriter,
	relativePath dir_layout.RelativePath,
) *BoxCheckedOut {
	return &BoxCheckedOut{
		BoxTransacted: BoxTransacted{
			ColorOptions:     co,
			Options:          options,
			Box:              fieldsFormatWriter,
			Abbr:             abbr,
			FSItemReadWriter: fsItemReadWriter,
			RelativePath:     relativePath,
		},
	}
}

type BoxCheckedOut struct {
	BoxTransacted
}
