package box_format

import (
	"slices"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"code.linenisgreat.com/zit/go/zit/src/charlie/options_print"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_metadata_fmt"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func MakeBoxTransacted(
	co string_format_writer.ColorOptions,
	options options_print.V0,
	fieldsFormatWriter interfaces.StringFormatWriter[string_format_writer.Box],
	abbr ids.Abbr,
	fsItemReadWriter sku.FSItemReadWriter,
	relativePath dir_layout.RelativePath,
) *BoxTransacted {
	return &BoxTransacted{
		ColorOptions:     co,
		Options:          options,
		Box:              fieldsFormatWriter,
		Abbr:             abbr,
		FSItemReadWriter: fsItemReadWriter,
		RelativePath:     relativePath,
	}
}

type BoxTransacted struct {
	string_format_writer.ColorOptions
	Options options_print.V0

	MaxHead, MaxTail int
	Padding          string

	Box interfaces.StringFormatWriter[string_format_writer.Box]

	ids.Abbr
	FSItemReadWriter sku.FSItemReadWriter
	dir_layout.RelativePath
}

func (f *BoxTransacted) WriteStringFormat(
	sw interfaces.WriterAndStringWriter,
	sk *sku.Transacted,
) (n int64, err error) {
	var box string_format_writer.Box

	var n2 int64

	box.Header.RightAligned = true

	if f.Options.PrintTime && !f.Options.PrintTai {
		t := sk.GetTai()
		box.Header.Value = t.Format(string_format_writer.StringFormatDateTime)
	}

	box.Contents = slices.Grow(box.Contents, 10)

	n2, err = f.addFieldsExternal2(
		sk,
		&box,
		f.Options.DescriptionInBox,
	)
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	b := &sk.Metadata.Description

	if !f.Options.DescriptionInBox && !b.IsEmpty() {
		box.Trailer = append(
			box.Trailer,
			string_format_writer.Field{
				Value:              b.StringWithoutNewlines(),
				ColorType:          string_format_writer.ColorTypeUserData,
				DisableValueQuotes: true,
			},
		)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	n2, err = f.Box.WriteStringFormat(sw, box)
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *BoxTransacted) addFieldsExternal(
	e *sku.Transacted,
	box *string_format_writer.Box,
	includeDescriptionInBox bool,
) (n int64, err error) {
	if e.State == external_state.Unknown {
		if e.ObjectId.IsEmpty() {
			e.State = external_state.Untracked
		}
	}

	oid := &e.ObjectId
	oidIsExternal := e.State == external_state.Untracked && !e.ExternalObjectId.IsEmpty()

	if oidIsExternal {
		oid = &e.ExternalObjectId
	}

	oidString := (*ids.ObjectIdStringerSansRepo)(oid).String()

	if f.Abbr.ZettelId.Abbreviate != nil &&
		oid.GetGenre() == genres.Zettel &&
		!oidIsExternal {
		if oidString, err = f.Abbr.ZettelId.Abbreviate(oid); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	objectIDField := string_format_writer.Field{
		Value: oidString,
		// DisableValueQuotes: oid.GetGenre() != genres.Blob,
		// DisableValueQuotes: true,
		ColorType: string_format_writer.ColorTypeId,
	}

	box.Contents = append(box.Contents, objectIDField)

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

	if err = f.addFieldsMetadata(
		o,
		e,
		includeDescriptionInBox,
		box,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *BoxTransacted) addFieldsMetadata(
	options options_print.V0,
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

	if options.PrintTai && o.GetGenre() != genres.InventoryList {
		box.Contents = append(
			box.Contents,
			object_metadata_fmt.MetadataFieldTai(m),
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

	return
}
