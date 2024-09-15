package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/erworben_cli_print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
)

type CliExternal struct {
	options erworben_cli_print_options.PrintOptions

	rightAlignedWriter         interfaces.StringFormatWriter[string]
	shaStringFormatWriter      interfaces.StringFormatWriter[interfaces.Sha]
	objectIdStringFormatWriter interfaces.StringFormatWriter[*ids.ObjectId]
	fdStringFormatWriter       interfaces.StringFormatWriter[*fd.FD]
	metadataStringFormatWriter interfaces.StringFormatWriter[*object_metadata.Metadata]

	transactedWriter *sku_fmt.Box
}

func MakeCliExternalFormat(
	options erworben_cli_print_options.PrintOptions,
	shaStringFormatWriter interfaces.StringFormatWriter[interfaces.Sha],
	fdStringFormatWriter interfaces.StringFormatWriter[*fd.FD],
	objectIdStringFormatWriter interfaces.StringFormatWriter[*ids.ObjectId],
	metadataStringFormatWriter interfaces.StringFormatWriter[*object_metadata.Metadata],
) *CliExternal {
	return &CliExternal{
		options:                    options,
		rightAlignedWriter:         string_format_writer.MakeRightAligned(),
		shaStringFormatWriter:      shaStringFormatWriter,
		objectIdStringFormatWriter: objectIdStringFormatWriter,
		fdStringFormatWriter:       fdStringFormatWriter,
		metadataStringFormatWriter: metadataStringFormatWriter,
	}
}

func (f *CliExternal) WriteStringFormat(
	sw interfaces.WriterAndStringWriter,
	col sku.ExternalLike,
) (n int64, err error) {
	var co *External

	switch colt := col.(type) {
	case *sku.Transacted:
		if n, err = f.transactedWriter.WriteStringFormat(sw, col); err != nil {
			err = errors.Wrap(err)
			return
		}

		return

	case *External:
		co = colt

	default:
		err = errors.Errorf("unsupported ExternalLike: %T", col)
		return
	}

	var (
		n1 int
		n2 int64
	)

	o := co
	fds := o.GetFDs()
	n1, err = sw.WriteString("[")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	var m checkout_mode.Mode

	if m, err = fds.GetCheckoutModeOrError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	firstFD := &fds.Object

	if m == checkout_mode.BlobOnly || m == checkout_mode.BlobRecognized {
		firstFD = &fds.Blob
	}

	n2, err = f.fdStringFormatWriter.WriteStringFormat(sw, firstFD)
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n2, err = f.metadataStringFormatWriter.WriteStringFormat(
    sw,
    o.Transacted.GetMetadata(),
  )
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if m == checkout_mode.BlobRecognized ||
		(m != checkout_mode.MetadataOnly && m != checkout_mode.None) {
		n2, err = f.writeStringFormatBlobFDsExcept(sw, fds, firstFD)
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

func (f *CliExternal) writeStringFormatBlobFDsExcept(
	sw interfaces.WriterAndStringWriter,
	fds *FDSet,
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

func (f *CliExternal) writeStringFormatBlobFD(
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

func (f *CliExternal) writeStringFormatUntracked(
	sw interfaces.WriterAndStringWriter,
	co *External,
	mode checkout_mode.Mode,
) (n int64, err error) {
	var (
		n1 int
		n2 int64
	)

	o := co
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

	n2, err = f.metadataStringFormatWriter.WriteStringFormat(
    sw,
    o.Transacted.GetMetadata(),
  )
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n2, err = f.writeStringFormatBlobFDsExcept(sw, fds, fdToPrint)
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
