package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/erworben_cli_print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type cliCheckedOut struct {
	options erworben_cli_print_options.PrintOptions

	rightAlignedWriter          interfaces.StringFormatWriter[string]
	shaStringFormatWriter       interfaces.StringFormatWriter[interfaces.ShaLike]
	kennungStringFormatWriter   interfaces.StringFormatWriter[*kennung.Id]
	fdStringFormatWriter        interfaces.StringFormatWriter[*fd.FD]
	metadateiStringFormatWriter interfaces.StringFormatWriter[*metadatei.Metadatei]
}

func MakeCliCheckedOutFormat(
	options erworben_cli_print_options.PrintOptions,
	shaStringFormatWriter interfaces.StringFormatWriter[interfaces.ShaLike],
	fdStringFormatWriter interfaces.StringFormatWriter[*fd.FD],
	kennungStringFormatWriter interfaces.StringFormatWriter[*kennung.Id],
	metadateiStringFormatWriter interfaces.StringFormatWriter[*metadatei.Metadatei],
) *cliCheckedOut {
	return &cliCheckedOut{
		options:                     options,
		rightAlignedWriter:          string_format_writer.MakeRightAligned(),
		shaStringFormatWriter:       shaStringFormatWriter,
		kennungStringFormatWriter:   kennungStringFormatWriter,
		fdStringFormatWriter:        fdStringFormatWriter,
		metadateiStringFormatWriter: metadateiStringFormatWriter,
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

		if co.State == checked_out_state.StateError {
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

	o := &co.External
	fds := o.GetFDsPtr()
	n1, err = sw.WriteString("[")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	m := fds.GetCheckoutMode()

	switch {
	case co.IsImport:
		fallthrough

	case m == checkout_mode.ModeAkteOnly:
		n2, err = f.kennungStringFormatWriter.WriteStringFormat(sw, &o.Kennung)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		n2, err = f.fdStringFormatWriter.WriteStringFormat(sw, &fds.Objekte)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if co.State == checked_out_state.StateConflicted {
		n1, err = sw.WriteString("]")
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	n2, err = f.metadateiStringFormatWriter.WriteStringFormat(sw, o.GetMetadatei())
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if m != checkout_mode.ModeObjekteOnly && m != checkout_mode.ModeNone {
		n1, err = sw.WriteString("\n")
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n2, err = f.rightAlignedWriter.WriteStringFormat(
			sw,
			"",
		)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n1, err = sw.WriteString(" ")
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n2, err = f.fdStringFormatWriter.WriteStringFormat(
			sw,
			&fds.Akte,
		)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	n1, err = sw.WriteString("]")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
