package sku_fmt

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (f *Box) WriteStringFormatFSBox(
	sw interfaces.WriterAndStringWriter,
	co *sku.CheckedOut,
	o *sku.Transacted,
  box *string_format_writer.Box,
	fds *sku.FSItem,
) (n int64, err error) {
	var n2 int64

	var m checkout_mode.Mode

	if m, err = fds.GetCheckoutModeOrError(); err != nil {
		if co.State == checked_out_state.Conflicted {
			err = nil
			m = checkout_mode.BlobOnly
		} else {
			err = errors.Wrapf(err, "FDs: %s", fds.Debug())
			return
		}
	}

	var fdAlreadyWritten *fd.FD

	op := f.Options
	op.ExcludeFields = true

	switch {
	case co.State == checked_out_state.Untracked:
		if err = f.writeStringFormatUntracked(co, m, box); err != nil {
			err = errors.Wrap(err)
			return
		}

	case co.IsImport:
		fallthrough

	case m == checkout_mode.BlobOnly || m == checkout_mode.BlobRecognized:
		box.Contents = append(
			box.Contents,
			string_format_writer.Field{
				Value:              (*ids.ObjectIdStringerSansRepo)(&o.ObjectId).String(),
				DisableValueQuotes: true,
				ColorType:          string_format_writer.ColorTypeId,
			},
		)

	default:
		box.Contents = append(
			box.Contents,
			string_format_writer.Field{
				ColorType: string_format_writer.ColorTypeId,
				Value:     f.Rel(fds.Object.GetPath()),
			},
		)

		fdAlreadyWritten = &fds.Object
	}

	if co.State != checked_out_state.Conflicted {
		if err = f.WriteMetadataToBox(
			op,
			o,
			true,
			box,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if m == checkout_mode.BlobRecognized ||
			(m != checkout_mode.MetadataOnly && m != checkout_mode.None) {
			if err = f.writeStringFormatBlobFDsExcept(
				fds,
				fdAlreadyWritten,
				box,
			); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	n2, err = f.Fields.WriteStringFormat(
		sw,
		*box,
	)
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *Box) writeStringFormatBlobFDsExcept(
	fds *sku.FSItem,
	except *fd.FD,
	box *string_format_writer.Box,
) (err error) {
	if fds.MutableSetLike == nil {
		err = errors.Errorf("FDSet.MutableSetLike was nil")
		return
	}

	for fd := range fds.All() {
		if except != nil && fd.Equals(except) {
			return
		}

		if err = f.writeStringFormatBlobFD(fd, box); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (f *Box) writeStringFormatBlobFD(
	fd *fd.FD,
	box *string_format_writer.Box,
) (err error) {
	if err != nil {
		err = errors.Wrap(err)
		return
	}

	box.Contents = append(
		box.Contents,
		string_format_writer.Field{
			NeedsNewline: true,
			Prefix:       string_format_writer.StringIndentWithSpace,
			Value:        f.Rel(fd.GetPath()),
			ColorType:    string_format_writer.ColorTypeId,
		},
	)

	return
}

func (f *Box) writeStringFormatUntracked(
	co *sku.CheckedOut,
	mode checkout_mode.Mode,
	box *string_format_writer.Box,
) (err error) {
	o := &co.External
	var i *sku.FSItem

	if i, err = f.FSItemReadWriter.ReadFSItemFromExternal(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	fdToPrint := &i.Blob

	if o.GetGenre() != genres.Zettel && !i.Object.IsEmpty() {
		fdToPrint = &i.Object
	}

	box.Contents = append(
		box.Contents,
		string_format_writer.Field{
			ColorType: string_format_writer.ColorTypeId,
			Value:     f.Rel(fdToPrint.GetPath()),
		},
	)

	if err = f.writeStringFormatBlobFDsExcept(i, fdToPrint, box); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
