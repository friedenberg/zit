package sku

import (
	"fmt"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/token_types"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/charlie/options_print"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/descriptions"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/echo/query_spec"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/id_fmts"
)

type ObjectIdAlignedFormat interface {
	SetMaxKopfUndSchwanz(kop, schwanz int)
}

func MakeFormatOrganize(
	options options_print.V0,
	shaStringFormatWriter interfaces.StringFormatWriter[interfaces.Sha],
	objectIdStringFormatWriter id_fmts.Aligned,
	typeStringFormatWriter interfaces.StringFormatWriter[*ids.Type],
	descriptionStringFormatWriter interfaces.StringFormatWriter[*descriptions.Description],
	tagsStringFormatWriter interfaces.StringFormatWriter[*ids.Tag],
) *Organize {
	options.PrintTime = false
	options.PrintShas = false

	return &Organize{
		options:                       options,
		shaStringFormatWriter:         shaStringFormatWriter,
		objectIdStringFormatWriter:    objectIdStringFormatWriter,
		typStringFormatWriter:         typeStringFormatWriter,
		descriptionStringFormatWriter: descriptionStringFormatWriter,
		etikettenStringFormatWriter:   tagsStringFormatWriter,
	}
}

type Organize struct {
	options options_print.V0

	maxKopf, maxSchwanz int
	padding             string

	shaStringFormatWriter         interfaces.StringFormatWriter[interfaces.Sha]
	objectIdStringFormatWriter    id_fmts.Aligned
	typStringFormatWriter         interfaces.StringFormatWriter[*ids.Type]
	descriptionStringFormatWriter interfaces.StringFormatWriter[*descriptions.Description]
	etikettenStringFormatWriter   interfaces.StringFormatWriter[*ids.Tag]
}

func (f *Organize) SetMaxKopfUndSchwanz(k, s int) {
	f.maxKopf, f.maxSchwanz = k, s
	f.padding = strings.Repeat(" ", 5+k+s)
	f.objectIdStringFormatWriter.SetMaxKopfUndSchwanz(k, s)
}

func (f *Organize) WriteStringFormat(
	sw interfaces.WriterAndStringWriter,
	o *Transacted,
) (n int64, err error) {
	var n1 int

	if f.options.PrintTime {
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
	n2, err = f.objectIdStringFormatWriter.WriteStringFormat(sw, &o.ObjectId)
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	sh := o.GetBlobSha()

	if f.options.PrintShas && (!sh.IsNull() || f.options.PrintEmptyShas) {
		n1, err = sw.WriteString("@")
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n2, err = f.shaStringFormatWriter.WriteStringFormat(sw, o.GetBlobSha())
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	t := o.GetMetadata().GetTypePtr()

	if len(t.String()) > 0 {
		if f.padding == "" {
			n1, err = sw.WriteString(" !")
		} else {
			n1, err = sw.WriteString("  !")
		}

		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n2, err = f.typStringFormatWriter.WriteStringFormat(sw, t)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	b := &o.Metadata.Description

	if f.options.PrintTagsAlways {
		b := o.GetMetadata().GetTags()

		for _, v := range quiter.SortedValues(b) {
			if f.options.ZittishNewlines {
				n1, err = fmt.Fprintf(sw, "\n%s", f.padding)
			} else {
				n1, err = sw.WriteString(" ")
			}

			n += int64(n1)

			if err != nil {
				err = errors.Wrap(err)
				return
			}

			n2, err = f.etikettenStringFormatWriter.WriteStringFormat(sw, &v)
			n += n2

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	if f.options.ZittishNewlines {
		n1, err = fmt.Fprintf(sw, "\n%s]", f.padding)
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

		n2, err = f.descriptionStringFormatWriter.WriteStringFormat(sw, b)
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
	o *Transacted,
) (n int64, err error) {
	if err = f.readStringFormatWithinBrackets(rb, o); err != nil {
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

	if err = o.Metadata.Description.TodoSetSlice(sl); err != nil {
		err = errors.Wrap(err)
		return
	}

	rb.AdvanceRead(sl.Len())

	return
}

func (f *Organize) readStringFormatWithinBrackets(
	rb *catgut.RingBuffer,
	o *Transacted,
) (err error) {
	var ts query_spec.TokenScanner
	ts.Reset(catgut.MakeRingBufferRuneScanner(rb))

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
			if err = o.ObjectId.TodoSetBytes(t); err != nil {
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
					err = errors.Wrapf(err, "Readable: %q", rb.PeekReadable())
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
					return
				}

				k.Reset()
			}

		default:
			err = errors.Errorf("invalid state: %d", state)
			return
		}
	}

	rb.AdvanceRead(n)

	return
}
