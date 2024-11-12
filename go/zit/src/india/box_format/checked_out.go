package box_format

import (
	"slices"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/options_print"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_metadata_fmt"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func MakeBoxCheckedOut(
	co string_format_writer.ColorOptions,
	options options_print.V0,
	fieldsFormatWriter interfaces.StringFormatWriter[string_format_writer.Box],
	abbr ids.Abbr,
	fsItemReadWriter sku.FSItemReadWriter,
	relativePath dir_layout.RelativePath,
	headerWriter string_format_writer.HeaderWriter[*sku.CheckedOut],
) *BoxCheckedOut {
	return &BoxCheckedOut{
		headerWriter: headerWriter,
		BoxTransacted: BoxTransacted{
			optionsColor:     co,
			optionsPrint:     options,
			box:              fieldsFormatWriter,
			abbr:             abbr,
			fsItemReadWriter: fsItemReadWriter,
			relativePath:     relativePath,
		},
	}
}

type BoxCheckedOut struct {
	headerWriter string_format_writer.HeaderWriter[*sku.CheckedOut]
	BoxTransacted
}

func (f *BoxCheckedOut) WriteStringFormat(
	sw interfaces.WriterAndStringWriter,
	co *sku.CheckedOut,
) (n int64, err error) {
	var box string_format_writer.Box

	if f.headerWriter != nil {
		if err = f.headerWriter.WriteBoxHeader(&box.Header, co); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	box.Contents = slices.Grow(box.Contents, 10)

	var fds *sku.FSItem
	var errFS error

	external := co.GetSkuExternal()

	if f.fsItemReadWriter != nil {
		fds, errFS = f.fsItemReadWriter.ReadFSItemFromExternal(external)
	}

	if f.fsItemReadWriter == nil || errFS != nil || !external.RepoId.IsEmpty() {
		if err = f.addFieldsExternalWithFSItem(
			external,
			&box,
			f.optionsPrint.DescriptionInBox,
			fds,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = f.addFieldsFS(co, &box, fds); err != nil {
			err = errors.Wrapf(err, "CheckedOut: %s", co)
			return
		}
	}

	b := &external.Metadata.Description

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

func (f *BoxCheckedOut) addFieldsExternalWithFSItem(
	external *sku.Transacted,
	box *string_format_writer.Box,
	includeDescriptionInBox bool,
	item *sku.FSItem,
) (err error) {
	if err = f.addFieldsObjectIdsWithFSItem(
		external,
		box,
		item,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = f.addFieldsMetadataWithFSItem(
		f.optionsPrint,
		external,
		includeDescriptionInBox,
		box,
		item,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *BoxCheckedOut) makeFieldExternalObjectIdsIfNecessary(
	sk *sku.Transacted,
	item *sku.FSItem,
) (field string_format_writer.Field, err error) {
	field = string_format_writer.Field{
		ColorType: string_format_writer.ColorTypeId,
	}

	switch {
	case (item == nil || item.Len() == 0) && !sk.ExternalObjectId.IsEmpty():
		oid := &sk.ExternalObjectId
		// TODO quote as necessary
		field.Value = (*ids.ObjectIdStringerSansRepo)(oid).String()

	case item != nil:
		switch item.GetCheckoutMode() {
		case checkout_mode.MetadataOnly, checkout_mode.MetadataAndBlob:
			// TODO quote as necessary
			if !item.Object.IsStdin() {
				field.Value = f.relativePath.Rel(item.Object.String())
			}
		}
	}

	return
}

func (f *BoxCheckedOut) makeFieldObjectId(
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

func (f *BoxCheckedOut) addFieldsObjectIdsWithFSItem(
	sk *sku.Transacted,
	box *string_format_writer.Box,
	item *sku.FSItem,
) (err error) {
	var external string_format_writer.Field

	if external, err = f.makeFieldExternalObjectIdsIfNecessary(
		sk,
		item,
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

func (f *BoxCheckedOut) addFieldsMetadataWithFSItem(
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

	if !options.ExcludeFields && (item == nil || item.Len() == 0) {
		box.Contents = append(box.Contents, m.Fields...)
	}

	return
}

func (f *BoxCheckedOut) addFieldsFS(
	co *sku.CheckedOut,
	box *string_format_writer.Box,
	item *sku.FSItem,
) (err error) {
	m := item.GetCheckoutMode()

	var fdAlreadyWritten *fd.FD

	op := f.optionsPrint
	op.ExcludeFields = true

	if co.GetState() == checked_out_state.Unknown {
		err = errors.Errorf("invalid state unknown")
		return
	}

	if co.GetState() == checked_out_state.Untracked {
		if err = f.addFieldsUntracked(co, box, item, op); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	// if co.State == checked_out_state.Untracked ||
	// 	co.State == checked_out_state.Unknown {
	// 	if err = f.addFieldsUntracked(co, box, item, op); err != nil {
	// 		err = errors.Wrap(err)
	// 		return
	// 	}

	// 	return
	// }

	var id string_format_writer.Field

	switch {
	// case co.IsImport:
	// 	fallthrough

	case m == checkout_mode.BlobOnly || m == checkout_mode.BlobRecognized:
		id.Value = (*ids.ObjectIdStringerSansRepo)(&co.GetSkuExternal().ObjectId).String()

	case m.IncludesMetadata():
		id.Value = f.relativePath.Rel(item.Object.GetPath())
		fdAlreadyWritten = &item.Object

	default:
		if id, _, err = f.makeFieldObjectId(
			co.GetSkuExternal(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	id.ColorType = string_format_writer.ColorTypeId
	box.Contents = append(box.Contents, id)

	if co.GetState() == checked_out_state.Conflicted {
	} else {
		if err = f.addFieldsMetadata(
			op,
			co.GetSkuExternal(),
			op.DescriptionInBox,
			box,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if m == checkout_mode.BlobRecognized ||
			(m != checkout_mode.MetadataOnly && m != checkout_mode.None) {
			if err = f.addFieldsFSBlobExcept(
				item,
				fdAlreadyWritten,
				box,
			); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	return
}

func (f *BoxTransacted) addFieldsFSBlobExcept(
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
			continue
		}

		if err = f.addFieldFSBlob(fd, box); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (f *BoxTransacted) addFieldFSBlob(
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
			Value:        f.relativePath.Rel(fd.GetPath()),
			ColorType:    string_format_writer.ColorTypeId,
			NeedsNewline: true,
		},
	)

	return
}

func (f *BoxTransacted) addFieldsUntracked(
	co *sku.CheckedOut,
	box *string_format_writer.Box,
	item *sku.FSItem,
	op options_print.V0,
) (err error) {
	fdToPrint := &item.Blob

	if co.GetSkuExternal().GetGenre() != genres.Zettel && !item.Object.IsEmpty() {
		fdToPrint = &item.Object
	}

	box.Contents = append(
		box.Contents,
		string_format_writer.Field{
			ColorType: string_format_writer.ColorTypeId,
			Value:     f.relativePath.Rel(fdToPrint.GetPath()),
		},
	)

	if err = f.addFieldsMetadata(
		op,
		co.GetSkuExternal(),
		f.optionsPrint.DescriptionInBox,
		box,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
