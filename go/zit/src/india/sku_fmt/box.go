package sku_fmt

import (
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"code.linenisgreat.com/zit/go/zit/src/charlie/print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
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
		o,
	); errFS != nil || !isCO || !o.RepoId.IsEmpty() {
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
		n2, err = f.WriteStringFormatFSBox(sw, co, o, fds)
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

	if e.State == external_state.Untracked {
		oid = &e.ExternalObjectId
	}

	objectIDField := string_format_writer.Field{
		Value:              (*ids.ObjectIdStringerSansRepo)(oid).String(),
		DisableValueQuotes: true,
		ColorType:          string_format_writer.ColorTypeId,
	}

	fields = append(fields, objectIDField)

	var n2 int64

	if e.State != external_state.Untracked {
		if !e.ExternalObjectId.IsEmpty() && false {
			fields = append(
				fields,
				string_format_writer.Field{
					Value:              (*ids.ObjectIdStringerSansRepo)(&e.ExternalObjectId).String(),
					DisableValueQuotes: true,
					ColorType:          string_format_writer.ColorTypeId,
				},
			)
		}

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
