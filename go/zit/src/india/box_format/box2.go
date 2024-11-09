package box_format

import (
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"code.linenisgreat.com/zit/go/zit/src/charlie/options_print"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_metadata_fmt"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (f *Box) makeFieldExternalObjectIdsIfNecessary(
	sk *sku.Transacted,
	item *sku.FSItem,
) (field string_format_writer.Field, err error) {
	// TODO figure out why this is necessary and in this horribly place
	if sk.State == external_state.Unknown {
		if sk.ObjectId.IsEmpty() {
			sk.State = external_state.Untracked
		}
	}

	field = string_format_writer.Field{
		ColorType: string_format_writer.ColorTypeId,
	}

	switch {
	case item == nil && !sk.ExternalObjectId.IsEmpty():
		oid := &sk.ExternalObjectId
		// TODO quote as necessary
		field.Value = (*ids.ObjectIdStringerSansRepo)(oid).String()

	case item != nil:
		switch item.GetCheckoutMode() {
		case checkout_mode.MetadataOnly, checkout_mode.MetadataAndBlob:
			// TODO quote as necessary
			if !item.Object.IsStdin() {
				field.Value = f.RelativePath.Rel(item.Object.String())
			}
		}
	}

	return
}

func (f *Box) addFieldsObjectIds2(
	sk *sku.Transacted,
	box *string_format_writer.Box,
	item *sku.FSItem,
) (err error) {
	var eoidField string_format_writer.Field

	if eoidField, err = f.makeFieldExternalObjectIdsIfNecessary(
		sk,
		item,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var oidField string_format_writer.Field

	{
		oid := &sk.ObjectId

		if !oid.IsEmpty() {

			oidString := (*ids.ObjectIdStringerSansRepo)(oid).String()

			if f.Abbr.ZettelId.Abbreviate != nil &&
				oid.GetGenre() == genres.Zettel {
				if oidString, err = f.Abbr.ZettelId.Abbreviate(oid); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			oidField = string_format_writer.Field{
				Value:              oidString,
				DisableValueQuotes: true,
				ColorType:          string_format_writer.ColorTypeId,
			}
		}
	}

	switch {
	case oidField.Value != "" && eoidField.Value != "":
		if strings.HasPrefix(eoidField.Value, oidField.Value) {
			box.Contents = append(box.Contents, eoidField)
		} else {
			box.Contents = append(box.Contents, oidField, eoidField)
		}

	case oidField.Value != "":
		box.Contents = append(box.Contents, oidField)

	case eoidField.Value != "":
		box.Contents = append(box.Contents, eoidField)

	default:
		err = errors.Errorf("empty id")
		return
	}

	return
}

func (f *Box) addFieldsMetadata2(
	options options_print.V0,
	sk *sku.Transacted,
	includeDescriptionInBox bool,
	box *string_format_writer.Box,
	item *sku.FSItem,
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

	if !options.ExcludeFields && item == nil {
		box.Contents = append(box.Contents, m.Fields...)
	}

	return
}
