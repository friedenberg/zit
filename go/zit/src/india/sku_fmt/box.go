package sku_fmt

import (
	"slices"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"code.linenisgreat.com/zit/go/zit/src/charlie/print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_metadata_fmt"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func MakeBox(
	co string_format_writer.ColorOptions,
	options print_options.General,
	fieldsFormatWriter interfaces.StringFormatWriter[string_format_writer.Box],
	abbr ids.Abbr,
	fsItemReadWriter sku.FSItemReadWriter,
	relativePath fs_home.RelativePath,
) *Box {
	return &Box{
		ColorOptions:     co,
		Options:          options,
		Fields:           fieldsFormatWriter,
		Abbr:             abbr,
		FSItemReadWriter: fsItemReadWriter,
		RelativePath:     relativePath,
	}
}

type Box struct {
	string_format_writer.ColorOptions
	Options print_options.General

	MaxHead, MaxTail int
	Padding          string

	Fields interfaces.StringFormatWriter[string_format_writer.Box]

	ids.Abbr
	FSItemReadWriter sku.FSItemReadWriter
	fs_home.RelativePath
}

func (f *Box) SetMaxKopfUndSchwanz(k, s int) {
	f.MaxHead, f.MaxTail = k, s
	f.Padding = strings.Repeat(" ", 5+k+s)
}

func (f *Box) WriteStringFormat(
	sw interfaces.WriterAndStringWriter,
	el sku.ExternalLike,
) (n int64, err error) {
	o := el.GetSku()
	var box string_format_writer.Box

	var n2 int64

	co, isCO := el.(*sku.CheckedOut)

	var stateString string

	if isCO {
		state := co.GetState()

		if state == checked_out_state.Error {
			stateString = co.GetError().Error()
		} else {
			stateString = state.String()
		}

		o = &co.External
	} else if f.Options.PrintState {
		stateString = o.GetExternalState().String()
	}

	if stateString != "" {
		box.Header = append(
			box.Header,
			string_format_writer.Field{
				Value: stateString,
			},
		)
	} else if f.Options.PrintTime {
		t := o.GetTai()

		box.Header = append(
			box.Header,
			string_format_writer.Field{
				Value: t.Format(string_format_writer.StringFormatDateTime),
			},
		)
	}

	if fds, errFS := f.FSItemReadWriter.ReadFSItemFromExternal(
		o,
	); errFS != nil || !isCO || !o.RepoId.IsEmpty() {
		n2, err = f.WriteStringFormatExternal(
			sw,
			o,
			&box,
			f.Options.DescriptionInBox,
		)
		n += int64(n2)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	} else {
		n2, err = f.WriteStringFormatFSBox(sw, co, o, &box, fds)
		n += n2

		if err != nil {
			err = errors.Wrapf(err, "CheckedOut: %s", co)
			return
		}
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
	box *string_format_writer.Box,
	includeDescriptionInBox bool,
) (n int64, err error) {
	if e.State == external_state.Unknown {
		if e.ObjectId.IsEmpty() {
			e.State = external_state.Untracked
		}
	}

	box.Contents = slices.Grow(box.Contents, 10)

	oid := &e.ObjectId

	if e.State == external_state.Untracked && !e.ExternalObjectId.IsEmpty() {
		oid = &e.ExternalObjectId
	}

	objectIDField := string_format_writer.Field{
		Value:              (*ids.ObjectIdStringerSansRepo)(oid).String(),
		DisableValueQuotes: true,
		ColorType:          string_format_writer.ColorTypeId,
	}

	box.Contents = append(box.Contents, objectIDField)

	var n2 int64

	o := f.Options

	if e.State != external_state.Untracked {
		if !e.ExternalObjectId.IsEmpty() && false {
			box.Contents = append(
				box.Contents,
				string_format_writer.Field{
					Value:              (*ids.ObjectIdStringerSansRepo)(&e.ExternalObjectId).String(),
					DisableValueQuotes: true,
					ColorType:          string_format_writer.ColorTypeId,
				},
			)
		}
	}

	if err = f.WriteMetadataToBox(
		o,
		e,
		includeDescriptionInBox,
		box,
	); err != nil {
		err = errors.Wrap(err)
		return
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

func (f *Box) WriteMetadataToBox(
	options print_options.General,
	o *sku.Transacted,
	includeDescriptionInBox bool,
	box *string_format_writer.Box,
) (err error) {
	m := o.GetMetadata()

	if options.PrintShas &&
		(options.PrintEmptyShas || !m.Blob.IsNull()) {
		var shaString string

		if shaString, err = object_metadata_fmt.MetadataShaString(
			m,
			f.Abbr.Sha.Abbreviate,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		box.Contents = append(
			box.Contents,
			object_metadata_fmt.MetadataFieldShaString(shaString),
		)
	}

	if !m.Type.IsEmpty() {
		box.Contents = append(
			box.Contents,
			object_metadata_fmt.MetadataFieldType(m),
		)
	}

	b := m.Description

	if includeDescriptionInBox && !b.IsEmpty() {
		box.Contents = append(
			box.Contents,
			object_metadata_fmt.MetadataFieldDescription(m),
		)
	}

	box.Contents = append(
		box.Contents,
		object_metadata_fmt.MetadataFieldTags(m)...,
	)

	if !options.ExcludeFields {
		box.Contents = append(box.Contents, m.Fields...)
	}

	if !options.DescriptionInBox && !b.IsEmpty() {
		box.Trailer = []string_format_writer.Field{
			{
				Value:              b.String(),
				ColorType:          string_format_writer.ColorTypeUserData,
				DisableValueQuotes: true,
			},
		}
	}

	return
}
