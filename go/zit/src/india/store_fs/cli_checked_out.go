package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/erworben_cli_print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type cliCheckedOut struct {
	options erworben_cli_print_options.PrintOptions

	rightAlignedWriter          interfaces.StringFormatWriter[string]
	shaStringFormatWriter       interfaces.StringFormatWriter[interfaces.Sha]
	objectIdStringFormatWriter  interfaces.StringFormatWriter[*ids.ObjectId]
	fdStringFormatWriter        interfaces.StringFormatWriter[*fd.FD]
	metadateiStringFormatWriter interfaces.StringFormatWriter[*object_metadata.Metadata]
}

func MakeCliCheckedOutFormat(
	options erworben_cli_print_options.PrintOptions,
	shaStringFormatWriter interfaces.StringFormatWriter[interfaces.Sha],
	fdStringFormatWriter interfaces.StringFormatWriter[*fd.FD],
	kennungStringFormatWriter interfaces.StringFormatWriter[*ids.ObjectId],
	metadateiStringFormatWriter interfaces.StringFormatWriter[*object_metadata.Metadata],
) *cliCheckedOut {
	return &cliCheckedOut{
		options:                     options,
		rightAlignedWriter:          string_format_writer.MakeRightAligned(),
		shaStringFormatWriter:       shaStringFormatWriter,
		objectIdStringFormatWriter:  kennungStringFormatWriter,
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
	case co.State == checked_out_state.StateUntracked:
		n2, err = f.writeStringFormatUntracked(sw, co)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		return

	case co.IsImport:
		fallthrough

	case m == checkout_mode.ModeBlobOnly:
		n2, err = f.objectIdStringFormatWriter.WriteStringFormat(sw, &o.ObjectId)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		n2, err = f.fdStringFormatWriter.WriteStringFormat(sw, &fds.Object)
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

	n2, err = f.metadateiStringFormatWriter.WriteStringFormat(sw, o.GetMetadata())
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if m != checkout_mode.ModeMetadataOnly && m != checkout_mode.ModeNone {
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
			&fds.Blob,
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

func (f *cliCheckedOut) writeStringFormatUntracked(
	sw interfaces.WriterAndStringWriter,
	co *CheckedOut,
) (n int64, err error) {
	var (
		n1 int
		n2 int64
	)

	o := &co.External
	fds := o.GetFDsPtr()

	fdToPrint := &fds.Blob

	if o.GetGenre() != genres.Zettel {
		fdToPrint = &fds.Object
	}

	n2, err = f.fdStringFormatWriter.WriteStringFormat(
		sw,
		fdToPrint,
	)
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n2, err = f.metadateiStringFormatWriter.WriteStringFormat(sw, o.GetMetadata())
	n += n2

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
