package sku_fmt

import (
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"code.linenisgreat.com/zit/go/zit/src/charlie/print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/id_fmts"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_metadata_fmt"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func MakeBox(
	co string_format_writer.ColorOptions,
	options print_options.General,
	shaStringFormatWriter interfaces.StringFormatWriter[interfaces.Sha],
	objectIdStringFormatWriter id_fmts.Aligned,
	typeStringFormatWriter interfaces.StringFormatWriter[*ids.Type],
	tagsStringFormatWriter interfaces.StringFormatWriter[*ids.Tag],
	fieldsFormatWriter interfaces.StringFormatWriter[string_format_writer.Box],
	metadata interfaces.StringFormatWriter[*object_metadata.Metadata],
	abbr ids.Abbr,
	fsItemReadWriter sku.FSItemReadWriter,
	fdStringFormatWriter interfaces.StringFormatWriter[*fd.FD],
) *Box {
	return &Box{
		ColorOptions:         co,
		Options:              options,
		ShaString:            shaStringFormatWriter,
		ObjectId:             objectIdStringFormatWriter,
		Type:                 typeStringFormatWriter,
		TagString:            tagsStringFormatWriter,
		Fields:               fieldsFormatWriter,
		Metadata:             metadata,
		RightAligned:         string_format_writer.MakeRightAligned(),
		Abbr:                 abbr,
		FSItemReadWriter:     fsItemReadWriter,
		fdStringFormatWriter: fdStringFormatWriter,
	}
}

type Box struct {
	string_format_writer.ColorOptions
	Options print_options.General

	MaxHead, MaxTail int
	Padding          string

	RightAligned interfaces.StringFormatWriter[string]

	ShaString interfaces.StringFormatWriter[interfaces.Sha]
	ObjectId  id_fmts.Aligned
	Type      interfaces.StringFormatWriter[*ids.Type]
	TagString interfaces.StringFormatWriter[*ids.Tag]
	Fields    interfaces.StringFormatWriter[string_format_writer.Box]
	Metadata  interfaces.StringFormatWriter[*object_metadata.Metadata]

	ids.Abbr
	FSItemReadWriter     sku.FSItemReadWriter
	fdStringFormatWriter interfaces.StringFormatWriter[*fd.FD]
}

func (f *Box) SetMaxKopfUndSchwanz(k, s int) {
	f.MaxHead, f.MaxTail = k, s
	f.Padding = strings.Repeat(" ", 5+k+s)
	f.ObjectId.SetMaxKopfUndSchwanz(k, s)
}

func (f *Box) WriteStringFormat(
	sw interfaces.WriterAndStringWriter,
	el sku.ExternalLike,
) (n int64, err error) {
	o := el.GetSku()

	var n1 int
	var n2 int64

	objectForFDs := o

	co, isCO := el.(*sku.CheckedOut)

	if isCO {
		objectForFDs = &co.External
		state := co.GetState()
		var stateString string

		if state == checked_out_state.Error {
			stateString = co.GetError().Error()
		} else {
			stateString = state.String()
		}

		n2, err = f.RightAligned.WriteStringFormat(sw, stateString)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}

	} else if f.Options.PrintTime {
		t := o.GetTai()

		n1, err = sw.WriteString(t.Format(string_format_writer.StringFormatDateTime))
		n += int64(n1)

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
	}

	n1, err = sw.WriteString("[")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if fds, errFS := f.FSItemReadWriter.ReadFSItemFromExternal(
		objectForFDs,
	); errFS != nil || !isCO {
		n2, err = f.WriteStringFormatExternal(
			sw,
			o,
			f.Options.DescriptionInBox,
		)
		n += int64(n2)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		n2, err = f.WriteStringFormatFSBox(sw, co, objectForFDs, fds)
		n += n2

		if err != nil {
			err = errors.Wrapf(err, "CheckedOut: %s", co)
			return
		}
	}

	n1, err = sw.WriteString("]")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	b := &o.Metadata.Description

	if !f.Options.DescriptionInBox && !b.IsEmpty() {
		n2, err = f.Fields.WriteStringFormat(
			sw,
			string_format_writer.Box{
				Contents: []string_format_writer.Field{
					{
						Value:              b.String(),
						ColorType:          string_format_writer.ColorTypeUserData,
						DisableValueQuotes: true,
						Prefix:             " ",
					},
				},
			},
		)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (f *Box) WriteStringFormatExternal(
	sw interfaces.WriterAndStringWriter,
	e *sku.Transacted,
	includeDescriptionInBox bool,
) (n int64, err error) {
	fields := make([]string_format_writer.Field, 0, 10)

	oid := &e.ObjectId

	objectIDField := string_format_writer.Field{
		Value:              (*ids.ObjectIdStringerSansRepo)(oid).String(),
		DisableValueQuotes: true,
		ColorType:          string_format_writer.ColorTypeId,
	}

	fields = append(fields, objectIDField)

	var n2 int64

	if e.State != external_state.Untracked {
		m := &e.Metadata

		if f.Options.PrintShas {
			var shaString string

			if shaString, err = object_metadata_fmt.MetadataShaString(
				m,
				f.Abbr.Sha.Abbreviate,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			fields = append(
				fields,
				object_metadata_fmt.MetadataFieldShaString(shaString),
			)
		}

		if !m.Type.IsEmpty() {
			fields = append(
				fields,
				object_metadata_fmt.MetadataFieldType(m),
			)
		}

		if includeDescriptionInBox && !m.Description.IsEmpty() {
			fields = append(
				fields,
				object_metadata_fmt.MetadataFieldDescription(m),
			)
		}

		fields = append(
			fields,
			object_metadata_fmt.MetadataFieldTags(m)...,
		)
	}

	if !f.Options.ExcludeFields {
		fields = append(fields, e.Metadata.Fields...)
	}

	n2, err = f.Fields.WriteStringFormat(
		sw,
		string_format_writer.Box{Contents: fields},
	)
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *Box) WriteStringFormatExternalBoxUntracked(
	sw interfaces.WriterAndStringWriter,
	i *sku.Transacted,
	e *sku.Transacted,
	unboxedDescription bool,
) (n int64, err error) {
	if e.State != external_state.Untracked {
		err = errors.Errorf(
			"expected state %s but got %s",
			external_state.Untracked,
			e.State,
		)

		return
	}

	boxed := []string_format_writer.Field{}
	unboxed := []string_format_writer.Field{}

	var n2 int64

	n2, err = f.Fields.WriteStringFormat(
		sw,
		string_format_writer.Box{
			Contents: boxed,
			Box:      true,
			Trailer:  unboxed,
		},
	)
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

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
