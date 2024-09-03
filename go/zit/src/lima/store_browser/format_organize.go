package store_browser

import (
	"bytes"
	"fmt"
	"net/url"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
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
	f *sku_fmt.Organize,
) *Organize {
	return &Organize{
		Organize: f,
	}
}

type Organize struct {
	*sku_fmt.Organize
}

func (f *Organize) WriteStringFormat(
	sw interfaces.WriterAndStringWriter,
	el sku.ExternalLike,
) (n int64, err error) {
	o := el.GetSku()
	var n2 int64

	if e, hasNative := el.(*External); hasNative {
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
				[]string_format_writer.Field{
					{
						Value:              b.String(),
						ColorType:          string_format_writer.ColorTypeUserData,
						DisableValueQuotes: true,
						Prefix:             " ",
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
		n2, err = f.Organize.WriteStringFormat(sw, o)

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

	if err = f.readStringFormatWithinBrackets(&ts, el); err != nil {
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

	return
}

func (f *Organize) readStringFormatWithinBrackets(
	ts *query_spec.TokenScanner,
	el sku.ExternalLike,
) (err error) {
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
			if err = o.ObjectId.TodoSetBytesForgiving(t); err != nil {
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
				case query_spec.TokenTypeField, query_spec.TokenTypeLiteral:
					e, hasNative := el.(*External)

					if !hasNative {
						err = errors.Errorf("unsupported type: %T", el)
						return
					}

					left := string(tokenParts.Left)
					right := string(tokenParts.Right)

					switch left {
					case "id":
						if err = e.Item.SetId(right); err != nil {
							err = errors.Wrap(err)
							return
						}

						continue LOOP

					case "url":
						e.Item.Url = right
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

	idFieldValue := e.ObjectId.String()
	var n2 int64

	// TODO make this more robust
	switch e.State {
	case external_state.Tracked, external_state.Recognized:
		if i != nil {
			idFieldValue = i.ObjectId.String()
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
		fields = append(
			fields,
			object_metadata_fmt.MetadataFieldSha(m),
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

	var u *url.URL

	if u, err = item.GetUrl(); err != nil {
		err = errors.Wrap(err)
		return
	}

	fields = append(
		fields,
		string_format_writer.Field{
			Key:       "url",
			Value:     u.String(),
			ColorType: string_format_writer.ColorTypeUserData,
		},
	)

	n2, err = f.Fields.WriteStringFormat(sw, fields)
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
