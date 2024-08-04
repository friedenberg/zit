package browser

import (
	"fmt"
	"io"
	"net/url"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/descriptions"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/echo/query_spec"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
)

func MakeFormatOrganize(
	f *sku_fmt.Organize,
) *Organize {
	return &Organize{
		Organize: f,
		descriptionStringFormatWriter: descriptions.MakeCliFormat(
			descriptions.CliFormatTruncation66CharEllipsis,
			f.ColorOptions,
			true,
		),
	}
}

type Organize struct {
	*sku_fmt.Organize
	descriptionStringFormatWriter interfaces.StringFormatWriter[*descriptions.Description]
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

	e, hasNative := el.(*External)

	if hasNative {

		item := e.item
		browser := &e.browser

		{
			n1, err = sw.WriteString(" !")
			n += int64(n1)

			if err != nil {
				err = errors.Wrap(err)
				return
			}

			n2, err = f.TypeStringFormatWriter.WriteStringFormat(
				sw,
				&browser.Metadata.Type,
			)
			n += n2

			if err != nil {
				err = errors.Wrap(err)
				return
			}

			if !browser.Metadata.Description.IsEmpty() {
				n1, err = sw.WriteString(" ")
				n += int64(n1)

				if err != nil {
					err = errors.Wrap(err)
					return
				}

				n2, err = f.descriptionStringFormatWriter.WriteStringFormat(
					sw,
					&browser.Metadata.Description,
				)
				n += n2

				if err != nil {
					err = errors.Wrap(err)
					return
				}
			}
		}

		{
			n1, err = sw.WriteString("\n")
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

			var u *url.URL

			if u, err = item.GetUrl(); err != nil {
				err = errors.Wrap(err)
				return
			}

			n1, err = fmt.Fprintf(sw, "url=%q", u)
			n += int64(n1)

			if err != nil {
				err = errors.Wrap(err)
				return
			}

			n1, err = sw.WriteString("\n")
			n += int64(n1)

			if err != nil {
				err = errors.Wrap(err)
				return
			}

			for _, v := range iter.SortedValues(browser.Metadata.GetTags()) {
				n1, err = sw.WriteString(" ")
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
	if err = f.readStringFormatWithinBrackets(rb, el); err != nil {
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
	rb *catgut.RingBuffer,
	el sku.ExternalLike,
) (err error) {
	o := el.GetSku()

	rr := catgut.MakeRingBufferRuneScanner(rb)

	state := 0
	var k ids.ObjectId
	var t catgut.String
	var eof bool
	var n int

LOOP:
	for !eof {
		t.Reset()
		err = query_spec.NextToken(rr, &t)

		if err == io.EOF {
			err = nil
			eof = true
		} else if err != nil {
			err = errors.Wrap(err)
			return
		}

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
			if err = o.ObjectId.TodoSetBytes(&t); err != nil {
				o.ObjectId.Reset()
				return
			}

			state++

		case 2:
			if t.EqualsString("]") {
				break LOOP
			} else {
				if err = k.TodoSetBytes(&t); err != nil {
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
