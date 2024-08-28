package store_browser

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/erworben_cli_print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type cliCheckedOut struct {
	options erworben_cli_print_options.PrintOptions

	rightAlignedWriter   interfaces.StringFormatWriter[string]
	stringFormatExternal *Organize
}

func MakeCliCheckedOutFormat(
	options erworben_cli_print_options.PrintOptions,
	stringFormatExternal *Organize,
) *cliCheckedOut {
	return &cliCheckedOut{
		options:              options,
		rightAlignedWriter:   string_format_writer.MakeRightAligned(),
		stringFormatExternal: stringFormatExternal,
	}
}

func (f *cliCheckedOut) WriteStringFormat(
	sw interfaces.WriterAndStringWriter,
	col sku.CheckedOutLike,
) (n int64, err error) {
	co := col.(*CheckedOut)
	var (
		n1 int
		n2 int64
	)

	{
		var stateString string

		if co.State == checked_out_state.Error {
			stateString = co.Error.Error()
		} else {
			stateString = co.State.String()
		}

		n2, err = f.rightAlignedWriter.WriteStringFormat(sw, stateString)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	n1, err = sw.WriteString("[")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if co.State != checked_out_state.Untracked {
	}

	n2, err = f.stringFormatExternal.writeStringFormatExternal(sw, &co.External)
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = sw.WriteString("]")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
