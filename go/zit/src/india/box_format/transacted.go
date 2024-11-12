package box_format

import (
	"slices"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
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
	headerWriter string_format_writer.HeaderWriter[*sku.Transacted],
) *BoxTransacted {
	return &BoxTransacted{
		optionsColor:     co,
		optionsPrint:     options,
		box:              fieldsFormatWriter,
		abbr:             abbr,
		fsItemReadWriter: fsItemReadWriter,
		relativePath:     relativePath,
		headerWriter:     headerWriter,
	}
}

type BoxTransacted struct {
	optionsColor string_format_writer.ColorOptions
	optionsPrint options_print.V0

	box          interfaces.StringFormatWriter[string_format_writer.Box]
	headerWriter string_format_writer.HeaderWriter[*sku.Transacted]

	abbr             ids.Abbr
	fsItemReadWriter sku.FSItemReadWriter
	relativePath     dir_layout.RelativePath
}

func (f *BoxTransacted) WriteStringFormat(
	sw interfaces.WriterAndStringWriter,
	sk *sku.Transacted,
) (n int64, err error) {
	var box string_format_writer.Box

	// box.Header.RightAligned = true

	// if f.optionsPrint.PrintTime && !f.optionsPrint.PrintTai {
	// 	t := sk.GetTai()
	// 	box.Header.Value = t.Format(string_format_writer.StringFormatDateTime)
	// }

	if f.headerWriter != nil {
		if err = f.headerWriter.WriteBoxHeader(&box.Header, sk); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	box.Contents = slices.Grow(box.Contents, 10)

	if err = f.addFieldsObjectIds(
		sk,
		&box,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = f.addFieldsMetadata(
		f.optionsPrint,
		sk,
		f.optionsPrint.DescriptionInBox,
		&box,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	b := &sk.Metadata.Description

	if !f.optionsPrint.DescriptionInBox && !b.IsEmpty() {
		box.Trailer = append(
			box.Trailer,
			string_format_writer.Field{
				Value:              b.StringWithoutNewlines(),
				ColorType:          string_format_writer.ColorTypeUserData,
				DisableValueQuotes: true,
			},
		)
	}

	if n, err = f.box.WriteStringFormat(sw, box); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *BoxTransacted) makeFieldExternalObjectIdsIfNecessary(
	sk *sku.Transacted,
) (field string_format_writer.Field, err error) {
	field = string_format_writer.Field{
		ColorType: string_format_writer.ColorTypeId,
	}

	if !sk.ExternalObjectId.IsEmpty() {
		oid := &sk.ExternalObjectId
		// TODO quote as necessary
		field.Value = (*ids.ObjectIdStringerSansRepo)(oid).String()
	}

	return
}

func (f *BoxTransacted) makeFieldObjectId(
	sk *sku.Transacted,
) (field string_format_writer.Field, empty bool, err error) {
	oid := &sk.ObjectId

	empty = oid.IsEmpty()

	oidString := (*ids.ObjectIdStringerSansRepo)(oid).String()

	if f.abbr.ZettelId.Abbreviate != nil &&
		oid.GetGenre() == genres.Zettel {
		if oidString, err = f.abbr.ZettelId.Abbreviate(oid); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	field = string_format_writer.Field{
		Value:              oidString,
		DisableValueQuotes: true,
		ColorType:          string_format_writer.ColorTypeId,
	}

	return
}

func (f *BoxTransacted) addFieldsObjectIds(
	sk *sku.Transacted,
	box *string_format_writer.Box,
) (err error) {
	var external string_format_writer.Field

	if external, err = f.makeFieldExternalObjectIdsIfNecessary(
		sk,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var internal string_format_writer.Field
	var externalEmpty bool

	if internal, externalEmpty, err = f.makeFieldObjectId(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	switch {
	// case internal.Value != "" && external.Value != "":
	// 	if strings.HasPrefix(external.Value, strings.TrimPrefix(internal.Value, "!")) {
	// 		box.Contents = append(box.Contents, external)
	// 	} else {
	// 		box.Contents = append(box.Contents, internal, external)
	// 	}

	case externalEmpty && external.Value != "":
		box.Contents = append(box.Contents, external)

	case internal.Value != "":
		box.Contents = append(box.Contents, internal)

	case external.Value != "":
		box.Contents = append(box.Contents, external)

	default:
		err = errors.Errorf("empty id")
		return
	}

	return
}

func (f *BoxTransacted) addFieldsMetadata(
	options options_print.V0,
	sk *sku.Transacted,
	includeDescriptionInBox bool,
	box *string_format_writer.Box,
) (err error) {
	m := sk.GetMetadata()

	if options.PrintShas &&
		(options.PrintEmptyShas || !m.Blob.IsNull()) {
		var shaString string

		if shaString, err = object_metadata_fmt.MetadataShaString(
			m,
			f.abbr.Sha.Abbreviate,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		box.Contents = append(
			box.Contents,
			object_metadata_fmt.MetadataFieldShaString(shaString),
		)
	}

	if options.PrintTai && sk.GetGenre() != genres.InventoryList {
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
