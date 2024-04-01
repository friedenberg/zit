package sku_fmt

import (
	"io"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/src/charlie/erworben_cli_print_options"
	"code.linenisgreat.com/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/src/delta/zittish"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type KennungAlignedFormat interface {
	SetMaxKopfUndSchwanz(kop, schwanz int)
}

type Organize struct {
	maxKopf, maxSchwanz int
	ex                  kennung.Abbr
	options             erworben_cli_print_options.PrintOptions
}

func MakeOrganizeFormat(
	ex kennung.Abbr,
	options erworben_cli_print_options.PrintOptions,
) *Organize {
	return &Organize{
		ex:      ex,
		options: options,
	}
}

func (f *Organize) SetMaxKopfUndSchwanz(k, s int) {
	f.maxKopf = k
	f.maxSchwanz = s
}

func (f *Organize) WriteStringFormat(
	sw schnittstellen.WriterAndStringWriter,
	o *sku.Transacted,
) (n int64, err error) {
	var n1 int

	n1, err = sw.WriteString("[")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	h := kennung.Aligned(&o.Kennung, f.maxKopf, f.maxSchwanz)
	n1, err = sw.WriteString(h)
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = sw.WriteString("]")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if !o.Metadatei.Bezeichnung.IsEmpty() {
		n1, err = sw.WriteString(" ")
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n1, err = sw.WriteString(o.Metadatei.Description())
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (f *Organize) ReadStringFormat(
	rb *catgut.RingBuffer,
	o *sku.Transacted,
) (n int64, err error) {
	if err = f.readStringFormatWithKennung(rb, o); err != nil {
		if collections.IsErrNotFound(err) {
			err = nil
		} else {
			err = errors.WrapExcept(err, io.EOF, catgut.ErrBufferEmpty)
			return
		}
	}

	var sl catgut.Slice

	if sl, err = rb.PeekUptoAndIncluding('\n'); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = o.Metadatei.Bezeichnung.TodoSetSlice(sl); err != nil {
		err = errors.Wrap(err)
		return
	}

	rb.AdvanceRead(sl.Len())

	return
}

func (f *Organize) readStringFormatWithKennung(
	rb *catgut.RingBuffer,
	o *sku.Transacted,
) (err error) {
	state := 0
	n := 0
	eof := false
	var t catgut.String

	var sl catgut.Slice

	if sl, err = rb.PeekUptoAndIncluding(']'); err != nil {
		if !collections.IsErrNotFound(err) {
			err = errors.Wrap(err)
		}

		return
	}

	rr := catgut.MakeSliceRuneScanner(sl)

LOOP:
	for !eof {
		t.Reset()
		err = zittish.NextToken(rr, &t)

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
			if err = o.Kennung.TodoSetBytes(&t); err != nil {
				err = errors.Wrapf(err, "Readable: len: %d, cap %d, 1:%q, 2:%q", rb.Len(), rb.Cap(), string(rb.PeekReadable().First()), string(rb.PeekReadable().Second()))
				o.Kennung.Reset()
				return
			}

			state++

		case 2:
			if !t.EqualsString("]") {
				o.Kennung.Reset()
				return
			}

			if f.options.Abbreviations.Hinweisen {
				if err = f.ex.AbbreviateHinweisOnly(
					&o.Kennung,
				); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			break LOOP

		default:
			err = errors.Errorf("invalid state: %d", state)
			return
		}
	}

	rb.AdvanceRead(sl.Len())

	return
}
