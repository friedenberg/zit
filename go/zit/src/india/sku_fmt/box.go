package sku_fmt

import (
	"fmt"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/token_types"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/erworben_cli_print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/echo/query_spec"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/id_fmts"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type ObjectIdAlignedFormat interface {
	SetMaxKopfUndSchwanz(kop, schwanz int)
}

func MakeBox(
	co string_format_writer.ColorOptions,
	options erworben_cli_print_options.PrintOptions,
	shaStringFormatWriter interfaces.StringFormatWriter[interfaces.Sha],
	objectIdStringFormatWriter id_fmts.Aligned,
	typeStringFormatWriter interfaces.StringFormatWriter[*ids.Type],
	tagsStringFormatWriter interfaces.StringFormatWriter[*ids.Tag],
	fieldsFormatWriter interfaces.StringFormatWriter[string_format_writer.Fields],
	metadata interfaces.StringFormatWriter[*object_metadata.Metadata],
	abbr ids.Abbr,
) *Box {
	options.PrintTime = false
	options.PrintShas = false

	co.OffEntirely = true

	return &Box{
		ColorOptions: co,
		Options:      options,
		ShaString:    shaStringFormatWriter,
		ObjectId:     objectIdStringFormatWriter,
		Type:         typeStringFormatWriter,
		TagString:    tagsStringFormatWriter,
		Fields:       fieldsFormatWriter,
		Metadata:     metadata,
		RightAligned: string_format_writer.MakeRightAligned(),
		Abbr:         abbr,
	}
}

type Box struct {
	string_format_writer.ColorOptions
	Options erworben_cli_print_options.PrintOptions

	MaxHead, MaxTail int
	Padding          string

	RightAligned interfaces.StringFormatWriter[string]

	ShaString interfaces.StringFormatWriter[interfaces.Sha]
	ObjectId  id_fmts.Aligned
	Type      interfaces.StringFormatWriter[*ids.Type]
	TagString interfaces.StringFormatWriter[*ids.Tag]
	Fields    interfaces.StringFormatWriter[string_format_writer.Fields]
	Metadata  interfaces.StringFormatWriter[*object_metadata.Metadata]

	ids.Abbr
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

	if f.Options.PrintTime {
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

	var n2 int64
	n2, err = f.ObjectId.WriteStringFormat(sw, &o.ObjectId)
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	sh := o.GetBlobSha()

	if f.Options.PrintShas && (!sh.IsNull() || f.Options.PrintEmptyShas) {
		n1, err = sw.WriteString("@")
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n2, err = f.ShaString.WriteStringFormat(sw, o.GetBlobSha())
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	t := o.GetMetadata().GetTypePtr()

	if len(t.String()) > 0 {
		if f.Padding == "" {
			n1, err = sw.WriteString(" !")
		} else {
			n1, err = sw.WriteString("  !")
		}

		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n2, err = f.Type.WriteStringFormat(sw, t)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	b := &o.Metadata.Description

	if f.Options.PrintTagsAlways {
		b := o.GetMetadata().GetTags()

		for _, v := range iter.SortedValues(b) {
			if f.Options.ZittishNewlines {
				n1, err = fmt.Fprintf(sw, "\n%s", f.Padding)
			} else {
				n1, err = sw.WriteString(" ")
			}

			n += int64(n1)

			if err != nil {
				err = errors.Wrap(err)
				return
			}

			n2, err = f.TagString.WriteStringFormat(sw, &v)
			n += n2

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}
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

	if !b.IsEmpty() {
		n2, err = f.Fields.WriteStringFormat(
			sw,
			string_format_writer.Fields{
				Boxed: []string_format_writer.Field{
					{
						Value:              b.String(),
						DisableValueQuotes: true,
						ColorType:          string_format_writer.ColorTypeUserData,
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

func (f *Box) ReadStringFormat(
	rb *catgut.RingBuffer,
	el sku.ExternalLike,
) (n int64, err error) {
	var ts query_spec.TokenScanner
	ts.Reset(catgut.MakeRingBufferRuneScanner(rb))

	if err = f.readStringFormatBox(&ts, el); err != nil {
		err = errors.Wrap(err)
		return
	}

	n = ts.N()

	o := el.GetSku()

	if err = o.Metadata.Description.ReadFromRuneScanner(&ts); err != nil {
		err = errors.Wrap(err)
		return
	}

	n = ts.N()

	return
}

func (f *Box) readStringFormatBox(
	ts *query_spec.TokenScanner,
	el sku.ExternalLike,
) (err error) {
	o := el.GetSku()

	state := 0
	var k ids.ObjectId

LOOP:
	for ts.ScanIdentifierLikeSkipSpaces() {
		t, tokenType, tokenParts := ts.GetTokenAndTypeAndParts()

		if t.EqualsString(" ") || t.EqualsString("\n") {
			continue
		}

		switch state {
		case 0:
			if !t.EqualsString("[") {
				ts.Unscan()
				return
			}

			state++

		case 1:
			if t.Bytes()[0] == '/' {
				// TODO set external object ID
			} else if err = o.ObjectId.ReadFromToken(t); err != nil {
				o.ObjectId.Reset()
				return
			}

			state++

		case 2:
			if t.EqualsString("]") {
				break LOOP
			} else {
				if tokenType == token_types.TypeField {
					ui.Debug().Print(tokenParts)
					continue
				}

				if err = k.TodoSetBytes(t); err != nil {
					err = errors.Wrap(err)
					return
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
					err = errors.Wrapf(err, "Token: %q", t)
					return
				}

				k.Reset()
			}

		default:
			err = errors.Errorf("invalid state: %d", state)
			return
		}
	}

	if ts.Error() != nil {
		err = errors.Wrap(ts.Error())
		return
	}

	return
}