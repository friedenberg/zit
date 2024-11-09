package box_format

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/options_print"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_metadata_fmt"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (f *BoxTransacted) addFieldsExternal2(
	e *sku.Transacted,
	box *string_format_writer.Box,
	includeDescriptionInBox bool,
) (n int64, err error) {
	if err = f.addFieldsObjectIds2(
		e,
		box,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = f.addFieldsMetadata2(
		f.Options,
		e,
		includeDescriptionInBox,
		box,
	); err != nil {
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

	if f.Abbr.ZettelId.Abbreviate != nil &&
		oid.GetGenre() == genres.Zettel {
		if oidString, err = f.Abbr.ZettelId.Abbreviate(oid); err != nil {
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

func (f *BoxTransacted) addFieldsObjectIds2(
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

func (f *BoxTransacted) addFieldsMetadata2(
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
