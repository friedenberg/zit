package sku_fmt

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/erworben_cli_print_options"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/catgut"
	"github.com/friedenberg/zit/src/delta/zittish"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type KennungAlignedFormat interface {
	SetMaxKopfUndSchwanz(kop, schwanz int)
}

type organize struct {
	maxKopf, maxSchwanz int
	ex                  kennung.Abbr
	options             erworben_cli_print_options.PrintOptions
}

func MakeOrganizeFormat(
	ex kennung.Abbr,
	options erworben_cli_print_options.PrintOptions,
) *organize {
	return &organize{
		ex:      ex,
		options: options,
	}
}

func (f *organize) SetMaxKopfUndSchwanz(k, s int) {
	f.maxKopf = k
	f.maxSchwanz = s
}

func (f *organize) WriteStringFormat(
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

func (f *organize) ReadStringFormat(
	rb *catgut.RingBuffer,
	o *sku.Transacted,
) (n int64, err error) {
	if err = f.readStringFormatWithKennung(rb, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	var sl catgut.Slice
	sl, err = rb.PeekUptoAndIncluding('\n')

	if err != nil {
		return
	}

	if err = o.Metadatei.Bezeichnung.TodoSetSlice(sl); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err != nil {
		return
	}

	rb.AdvanceRead(sl.Len())

	return
}

func (f *organize) readStringFormatWithKennung(
	rb *catgut.RingBuffer,
	o *sku.Transacted,
) (err error) {
	state := 0
	n := 0
	eof := false
	var t catgut.String

	rr := catgut.MakeRingBufferRuneScanner(rb)

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

	rb.AdvanceRead(n)

	return
}
