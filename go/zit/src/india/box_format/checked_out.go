package box_format

import (
	"slices"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
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
) *BoxCheckedOut {
	return &BoxCheckedOut{
		BoxTransacted: BoxTransacted{
			ColorOptions:     co,
			Options:          options,
			Box:              fieldsFormatWriter,
			Abbr:             abbr,
			FSItemReadWriter: fsItemReadWriter,
			RelativePath:     relativePath,
		},
	}
}

type BoxCheckedOut struct {
	BoxTransacted
}

func (f *BoxCheckedOut) WriteStringFormat(
	sw interfaces.WriterAndStringWriter,
	co *sku.CheckedOut,
) (n int64, err error) {
	var box string_format_writer.Box

	var n2 int64

	var stateString string
	var isError bool

	state := co.GetState()

	if state == checked_out_state.Error {
		isError = true
	}

	stateString = state.String()
	o := &co.External

	box.Header.RightAligned = true

	if stateString != "" {
		box.Header.Value = stateString
	} else if f.Options.PrintTime && !f.Options.PrintTai {
		t := o.GetTai()
		box.Header.Value = t.Format(string_format_writer.StringFormatDateTime)
	}

	box.Contents = slices.Grow(box.Contents, 10)

	var fds *sku.FSItem
	var errFS error

	if f.FSItemReadWriter != nil {
		fds, errFS = f.FSItemReadWriter.ReadFSItemFromExternal(o)
	}

	if f.FSItemReadWriter == nil || errFS != nil || !o.RepoId.IsEmpty() || isError {
		n2, err = f.addFieldsExternalWithFSItem(
			o,
			&box,
			f.Options.DescriptionInBox,
			fds,
		)
		n += int64(n2)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		n2, err = f.addFieldsFS(co, o, &box, fds)
		n += n2

		if err != nil {
			err = errors.Wrapf(err, "CheckedOut: %s", co)
			return
		}
	}

	b := &o.Metadata.Description

	if !f.Options.DescriptionInBox && !b.IsEmpty() {
		box.Trailer = append(
			box.Trailer,
			string_format_writer.Field{
				Value:              b.StringWithoutNewlines(),
				ColorType:          string_format_writer.ColorTypeUserData,
				DisableValueQuotes: true,
			},
		)
	}

	n2, err = f.Box.WriteStringFormat(sw, box)
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *BoxCheckedOut) addFieldsExternal(
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

func (f *BoxCheckedOut) addFieldsMetadata(
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

func (f *BoxCheckedOut) addFieldsExternalWithFSItem(
	e *sku.Transacted,
	box *string_format_writer.Box,
	includeDescriptionInBox bool,
	item *sku.FSItem,
) (n int64, err error) {
	if err = f.addFieldsObjectIds2(
		e,
		box,
		item,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = f.addFieldsMetadata2(
		f.Options,
		e,
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
				field.Value = f.RelativePath.Rel(item.Object.String())
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

func (f *BoxCheckedOut) addFieldsObjectIds2(
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

func (f *BoxCheckedOut) addFieldsMetadata2(
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

	if !options.ExcludeFields && (item == nil || item.Len() == 0) {
		box.Contents = append(box.Contents, m.Fields...)
	}

	return
}

func (f *BoxCheckedOut) addFieldsFS(
	co *sku.CheckedOut,
	o *sku.Transacted,
	box *string_format_writer.Box,
	fds *sku.FSItem,
) (n int64, err error) {
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
		if err = f.addFieldsFSUntracked(co, m, box); err != nil {
			err = errors.Wrap(err)
			return
		}

	case co.IsImport:
		fallthrough

	case m == checkout_mode.BlobOnly || m == checkout_mode.BlobRecognized:
		box.Contents = append(
			box.Contents,
			string_format_writer.Field{
				Value: (*ids.ObjectIdStringerSansRepo)(&o.ObjectId).String(),
				// DisableValueQuotes: true,
				// ColorType: string_format_writer.ColorTypeId,
				// Value:     f.Rel(fds.Blob.GetPath()),
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
		if err = f.addFieldsMetadata(
			op,
			o,
			op.DescriptionInBox,
			box,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if m == checkout_mode.BlobRecognized ||
			(m != checkout_mode.MetadataOnly && m != checkout_mode.None) {
			if err = f.addFieldsFSBlobExcept(
				fds,
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
			Value:        f.Rel(fd.GetPath()),
			ColorType:    string_format_writer.ColorTypeId,
			NeedsNewline: true,
		},
	)

	return
}

func (f *BoxTransacted) addFieldsFSUntracked(
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

	if err = f.addFieldsFSBlobExcept(i, fdToPrint, box); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
