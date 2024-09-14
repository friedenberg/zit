package store_browser

import (
	"bytes"
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/token_types"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/echo/query_spec"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_metadata_fmt"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
)

func MakeFormatOrganize(
	f *sku_fmt.Box,
) *Organize {
	return &Organize{
		Box: f,
	}
}

type Organize struct {
	*sku_fmt.Box
}

func (f *Organize) WriteStringFormat(
	sw interfaces.WriterAndStringWriter,
	el sku.ExternalLike,
) (n int64, err error) {
	o := el.GetSku()
	var n2 int64

	if e, hasNative := el.(*External); hasNative {
		switch e.State {
		case external_state.Untracked:
			if n, err = f.writeStringFormatExternalBoxUntracked(
				sw,
				nil,
				e,
				true,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

		var n1 int
		n1, err = sw.WriteString("[")
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n2, err = f.writeStringFormatExternal(sw, nil, e, false)
		n += int64(n2)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		if f.Options.ZittishNewlines {
			n1, err = fmt.Fprintf(sw, "\n%s]", f.Padding)
		} else {
			n1, err = sw.WriteString("]")
		}

		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		b := &e.Transacted.Metadata.Description

		if !b.IsEmpty() {
			n2, err = f.Fields.WriteStringFormat(
				sw,
				string_format_writer.Fields{
					Boxed: []string_format_writer.Field{
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
	} else {
		n2, err = f.Box.WriteStringFormat(sw, o)

		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (f *Organize) ReadStringFormat(
	rb *catgut.RingBuffer,
	el sku.ExternalLike,
) (n int64, err error) {
	var ts query_spec.TokenScanner
	ts.Reset(catgut.MakeRingBufferRuneScanner(rb))

	if err = f.readStringFormatBoxContents(&ts, el); err != nil {
		err = errors.Wrap(err)
		return
	}

	var sl catgut.Slice

	if sl, err = rb.PeekUptoAndIncluding('\n'); err != nil {
		if collections.IsErrNotFound(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	o := el.GetSku()

	if err = o.Metadata.Description.TodoSetSlice(sl); err != nil {
		err = errors.Wrap(err)
		return
	}

	rb.AdvanceRead(sl.Len())

	if err = el.GetSkuExternalLike().GetSku().Metadata.Type.Set(
		"!toml-bookmark",
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *Organize) readStringFormatBoxContents(
	ts *query_spec.TokenScanner,
	el sku.ExternalLike,
) (err error) {
	e := el.(*External)
	o := el.GetSku()

	state := 0
	var k ids.ObjectId
	var n int

LOOP:
	for ts.Scan() {
		t, tokenType, tokenParts := ts.GetTokenAndTypeAndParts()
		n += t.Len()

		if t.EqualsString(" ") || t.EqualsString("\n") {
			continue
		}

		switch state {
		case 0:
			if !t.EqualsString("[") {
				return
			}

			state++

		case 1:
			if t.Bytes()[0] == '/' {
				if err = e.Item.Id.Set(t.String()); err != nil {
					err = errors.Wrap(err)
					return
				}
				// TODO set external Id to t
			} else if err = o.ObjectId.TodoSetBytesForgiving(t); err != nil {
				err = errors.Wrap(err)
				o.ObjectId.Reset()
				return
			}

			state++

		case 2:
			if t.EqualsString("]") {
				break LOOP
			} else {
				switch tokenType {
				case token_types.TypeField, token_types.TypeLiteral:
					e, hasNative := el.(*External)

					if !hasNative {
						err = errors.Errorf("unsupported type: %T", el)
						return
					}

					left := string(tokenParts.Left)
					right := string(tokenParts.Right)

					switch left {
					case "id":
						if err = e.Item.Id.Set(right); err != nil {
							err = errors.Wrap(err)
							return
						}

						continue LOOP

					case "url":
						if err = e.Item.Url.UnmarshalBinary(tokenParts.Right); err != nil {
							err = errors.Wrap(err)
							return
						}

						continue LOOP

					case "title":
						e.Item.Title = right
						continue LOOP

					case "":
						if err = e.Metadata.Description.Set(right); err != nil {
							err = errors.Wrap(err)
							return
						}

						continue LOOP

					default:
						err = errors.Errorf("unsupported field type: %q", tokenParts.Left)
						return
					}
				}

				if bytes.HasPrefix(t.Bytes(), []byte("@")) {
					continue LOOP
				} else {
					if err = k.TodoSetBytes(t); err != nil {
						err = errors.Wrap(err)
						return
					}
				}

				g := k.GetGenre()

				switch g {
				case genres.Type:
					if err = o.Metadata.Type.TodoSetFromObjectId(&k); err != nil {
						err = errors.Wrap(err)
						return
					}

				case genres.Tag:
					var e ids.Tag

					if err = e.TodoSetFromObjectId(&k); err != nil {
						err = errors.Wrap(err)
						return
					}

					if err = o.AddTagPtr(&e); err != nil {
						err = errors.Wrap(err)
						return
					}

				default:
					err = genres.MakeErrUnsupportedGenre(k.GetGenre())
					return
				}

				k.Reset()
			}

		default:
			err = errors.Errorf("invalid state: %d", state)
			return
		}
	}

	return
}

func (f *Organize) writeStringFormatExternal(
	sw interfaces.WriterAndStringWriter,
	i *sku.Transacted,
	e *External,
	includeDescriptionInBox bool,
) (n int64, err error) {
	fields := []string_format_writer.Field{}

	idFieldValue := (*ids.ObjectIdStringerSansRepo)(&e.ObjectId).String()
	var n2 int64

	// TODO make this more robust
	switch e.State {
	case external_state.Tracked, external_state.Recognized:
		if i != nil {
			idFieldValue = (*ids.ObjectIdStringerSansRepo)(&i.ObjectId).String()
		}

	case external_state.Untracked:
		idFieldValue = "/"
	}

	fields = append(
		fields,
		string_format_writer.Field{
			Value:              idFieldValue,
			DisableValueQuotes: true,
			ColorType:          string_format_writer.ColorTypeId,
		},
	)

	if e.State != external_state.Untracked {
		m := &e.Metadata

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
			object_metadata_fmt.MetadataFieldType(m),
		)

		if includeDescriptionInBox {
			fields = append(
				fields,
				object_metadata_fmt.MetadataFieldDescription(m),
			)
		}
	}

	item := &e.Item

	fields = append(
		fields,
		string_format_writer.Field{
			Key:       "id",
			Value:     item.Id.String(),
			ColorType: string_format_writer.ColorTypeId,
		},
	)

	if item.Title != "" {
		fields = append(fields,
			string_format_writer.Field{
				Key:       "title",
				Value:     item.Title,
				ColorType: string_format_writer.ColorTypeUserData,
			},
		)
	}

	fields = append(
		fields,
		string_format_writer.Field{
			Key:       "url",
			Value:     item.Url.String(),
			ColorType: string_format_writer.ColorTypeUserData,
		},
	)

	n2, err = f.Fields.WriteStringFormat(
		sw,
		string_format_writer.Fields{Boxed: fields},
	)
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *Organize) writeStringFormatExternalBoxUntracked(
	sw interfaces.WriterAndStringWriter,
	i *sku.Transacted,
	e *External,
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

	item := &e.Item

	boxed = append(
		boxed,
		string_format_writer.Field{
			Value:              item.Id.String(),
			DisableValueQuotes: true,
			ColorType:          string_format_writer.ColorTypeId,
		},
	)

	if item.Title != "" {
		field := string_format_writer.Field{
			Value:     item.Title,
			ColorType: string_format_writer.ColorTypeUserData,
		}

		if unboxedDescription {
			field.DisableValueQuotes = true
			unboxed = append(unboxed, field)
		} else {
			field.Key = "title"
			boxed = append(boxed, field)
		}
	}

	boxed = append(
		boxed,
		string_format_writer.Field{
			Key:       "url",
			Value:     item.Url.String(),
			ColorType: string_format_writer.ColorTypeUserData,
		},
	)

	n2, err = f.Fields.WriteStringFormat(
		sw,
		string_format_writer.Fields{
			Boxed:   boxed,
			Box:     true,
			Unboxed: unboxed,
		},
	)
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
