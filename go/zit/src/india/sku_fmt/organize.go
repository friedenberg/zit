package sku_fmt

import (
	"fmt"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/erworben_cli_print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/descriptions"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/echo/query_spec"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/id_fmts"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type ObjectIdAlignedFormat interface {
	SetMaxKopfUndSchwanz(kop, schwanz int)
}

func MakeFormatOrganize(
	co string_format_writer.ColorOptions,
	options erworben_cli_print_options.PrintOptions,
	shaStringFormatWriter interfaces.StringFormatWriter[interfaces.Sha],
	objectIdStringFormatWriter id_fmts.Aligned,
	typeStringFormatWriter interfaces.StringFormatWriter[*ids.Type],
	descriptionStringFormatWriter interfaces.StringFormatWriter[*descriptions.Description],
	tagsStringFormatWriter interfaces.StringFormatWriter[*ids.Tag],
) *Organize {
	options.PrintTime = false
	options.PrintShas = false

	return &Organize{
		ColorOptions:                  co,
		Options:                       options,
		ShaStringFormatWriter:         shaStringFormatWriter,
		ObjectIdStringFormatWriter:    objectIdStringFormatWriter,
		TypeStringFormatWriter:        typeStringFormatWriter,
		DescriptionStringFormatWriter: descriptionStringFormatWriter,
		TagStringFormatWriter:         tagsStringFormatWriter,
	}
}

type Organize struct {
	string_format_writer.ColorOptions
	Options erworben_cli_print_options.PrintOptions

	MaxHead, MaxTail int
	Padding          string

	ShaStringFormatWriter         interfaces.StringFormatWriter[interfaces.Sha]
	ObjectIdStringFormatWriter    id_fmts.Aligned
	TypeStringFormatWriter        interfaces.StringFormatWriter[*ids.Type]
	DescriptionStringFormatWriter interfaces.StringFormatWriter[*descriptions.Description]
	TagStringFormatWriter         interfaces.StringFormatWriter[*ids.Tag]
}

func (f *Organize) SetMaxKopfUndSchwanz(k, s int) {
	f.MaxHead, f.MaxTail = k, s
	f.Padding = strings.Repeat(" ", 5+k+s)
	f.ObjectIdStringFormatWriter.SetMaxKopfUndSchwanz(k, s)
}

func (f *Organize) WriteStringFormat(
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
	n2, err = f.ObjectIdStringFormatWriter.WriteStringFormat(sw, &o.ObjectId)
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

		n2, err = f.ShaStringFormatWriter.WriteStringFormat(sw, o.GetBlobSha())
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

		n2, err = f.TypeStringFormatWriter.WriteStringFormat(sw, t)
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

			n2, err = f.TagStringFormatWriter.WriteStringFormat(sw, &v)
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
		n1, err = sw.WriteString(" ")
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n2, err = f.DescriptionStringFormatWriter.WriteStringFormat(sw, b)
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

	n = ts.N()

	o := el.GetSku()

	if err = o.Metadata.Description.ReadFromTokenScanner(&ts); err != nil {
		err = errors.Wrap(err)
		return
	}

	n = ts.N()

	return
}

func (f *Organize) readStringFormatWithinBrackets(
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
			if err = o.ObjectId.ReadFromToken(t); err != nil {
				o.ObjectId.Reset()
				return
			}

			state++

		case 2:
			if t.EqualsString("]") {
				break LOOP
			} else {
				if tokenType == query_spec.TokenTypeField {
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
