package sku_fmt

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (f *Box) WriteStringFormatFSBox(
	sw interfaces.WriterAndStringWriter,
	co *sku.CheckedOut,
	o *sku.Transacted,
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
		n2, err = f.ObjectId.WriteStringFormat(
			sw,
			&o.ObjectId,
		)
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

		fdAlreadyWritten = &fds.Object
	}

	if co.State == checked_out_state.Conflicted {
		return
	}

	n2, err = f.Metadata.WriteStringFormat(
		sw,
		o.GetMetadata(),
	)
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if m == checkout_mode.BlobRecognized ||
		(m != checkout_mode.MetadataOnly && m != checkout_mode.None) {
		n2, err = f.writeStringFormatBlobFDsExcept(sw, fds, fdAlreadyWritten)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (f *Box) writeStringFormatBlobFDsExcept(
	sw interfaces.WriterAndStringWriter,
	fds *sku.FSItem,
	except *fd.FD,
) (n int64, err error) {
	var n2 int64

	if fds.MutableSetLike == nil {
		err = errors.Errorf("FDSet.MutableSetLike was nil")
		return
	}

	if err = fds.MutableSetLike.Each(
		func(fd *fd.FD) (err error) {
			if except != nil && fd.Equals(except) {
				return
			}

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

	return
}

func (f *Box) writeStringFormatBlobFD(
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

	n2, err = f.RightAligned.WriteStringFormat(
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

func (f *Box) writeStringFormatUntracked(
	sw interfaces.WriterAndStringWriter,
	co *sku.CheckedOut,
	mode checkout_mode.Mode,
) (n int64, err error) {
	var n2 int64

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

	n2, err = f.fdStringFormatWriter.WriteStringFormat(
		sw,
		fdToPrint,
	)
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n2, err = f.Metadata.WriteStringFormat(
		sw,
		o.GetMetadata(),
	)
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n2, err = f.writeStringFormatBlobFDsExcept(sw, i, fdToPrint)
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
