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

	rightAlignedWriter         interfaces.StringFormatWriter[string]
	shaStringFormatWriter      interfaces.StringFormatWriter[interfaces.Sha]
	objectIdStringFormatWriter interfaces.StringFormatWriter[*ids.ObjectId]
	fdStringFormatWriter       interfaces.StringFormatWriter[*fd.FD]
	metadataStringFormatWriter interfaces.StringFormatWriter[*object_metadata.Metadata]
}

func MakeCliCheckedOutFormat(
	options erworben_cli_print_options.PrintOptions,
	shaStringFormatWriter interfaces.StringFormatWriter[interfaces.Sha],
	fdStringFormatWriter interfaces.StringFormatWriter[*fd.FD],
	objectIdStringFormatWriter interfaces.StringFormatWriter[*ids.ObjectId],
	metadataStringFormatWriter interfaces.StringFormatWriter[*object_metadata.Metadata],
) *cliCheckedOut {
	return &cliCheckedOut{
		options:                    options,
		rightAlignedWriter:         string_format_writer.MakeRightAligned(),
		shaStringFormatWriter:      shaStringFormatWriter,
		objectIdStringFormatWriter: objectIdStringFormatWriter,
		fdStringFormatWriter:       fdStringFormatWriter,
		metadataStringFormatWriter: metadataStringFormatWriter,
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

	o := &co.External
	fds := o.GetFDs()
	n1, err = sw.WriteString("[")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	var m checkout_mode.Mode

	if m, err = fds.GetCheckoutModeOrError(); err != nil {
		if co.State == checked_out_state.Conflicted {
			err = nil
			m = checkout_mode.BlobOnly
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	switch {
	case co.State == checked_out_state.Untracked:
		n2, err = f.writeStringFormatUntracked(sw, co, m)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		return

	case co.IsImport:
		fallthrough

	case m == checkout_mode.BlobOnly || m == checkout_mode.BlobRecognized:
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

	if co.State == checked_out_state.Conflicted {
		n1, err = sw.WriteString("]")
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	n2, err = f.metadataStringFormatWriter.WriteStringFormat(sw, o.GetMetadata())
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if m == checkout_mode.BlobRecognized {
		if err = fds.MutableSetLike.Each(
			func(fd *fd.FD) (err error) {
				n2, err = f.writeStringFormatBlobFD(sw, fd)
				n += n2

				if err != nil {
					err = errors.Wrap(err)
					return
				}

				return
			},
		); err != nil {
			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	} else if m != checkout_mode.MetadataOnly && m != checkout_mode.None {
		n2, err = f.writeStringFormatBlobFD(sw, &fds.Blob)
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

func (f *cliCheckedOut) writeStringFormatBlobFD(
	sw interfaces.WriterAndStringWriter,
	fd *fd.FD,
) (n int64, err error) {
	var n1 int
	var n2 int64

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

	n2, err = f.fdStringFormatWriter.WriteStringFormat(sw, fd)
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *cliCheckedOut) writeStringFormatUntracked(
	sw interfaces.WriterAndStringWriter,
	co *CheckedOut,
	mode checkout_mode.Mode,
) (n int64, err error) {
	var (
		n1 int
		n2 int64
	)

	o := &co.External
	fds := o.GetFDs()

	fdToPrint := &fds.Blob

	if o.GetGenre() != genres.Zettel && !fds.Object.IsEmpty() {
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

	n2, err = f.metadataStringFormatWriter.WriteStringFormat(sw, o.GetMetadata())
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
