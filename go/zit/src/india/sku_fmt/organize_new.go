package sku_fmt

import (
	"fmt"
	"io"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/charlie/erworben_cli_print_options"
	"code.linenisgreat.com/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/src/echo/bezeichnung"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/echo/zittish"
	"code.linenisgreat.com/zit/src/foxtrot/kennung_fmt"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type OrganizeNew struct {
	options erworben_cli_print_options.PrintOptions

	shaStringFormatWriter         schnittstellen.StringFormatWriter[schnittstellen.ShaLike]
	kennungStringFormatWriter     kennung_fmt.Aligned
	typStringFormatWriter         schnittstellen.StringFormatWriter[*kennung.Typ]
	bezeichnungStringFormatWriter schnittstellen.StringFormatWriter[*bezeichnung.Bezeichnung]
	etikettenStringFormatWriter   schnittstellen.StringFormatWriter[*kennung.Etikett]
}

func MakeOrganizeNewFormat(
	options erworben_cli_print_options.PrintOptions,
	shaStringFormatWriter schnittstellen.StringFormatWriter[schnittstellen.ShaLike],
	kennungStringFormatWriter kennung_fmt.Aligned,
	typStringFormatWriter schnittstellen.StringFormatWriter[*kennung.Typ],
	bezeichnungStringFormatWriter schnittstellen.StringFormatWriter[*bezeichnung.Bezeichnung],
	etikettenStringFormatWriter schnittstellen.StringFormatWriter[*kennung.Etikett],
) *OrganizeNew {
	options.PrintTime = false
	options.PrintShas = false

	return &OrganizeNew{
		options:                       options,
		shaStringFormatWriter:         shaStringFormatWriter,
		kennungStringFormatWriter:     kennungStringFormatWriter,
		typStringFormatWriter:         typStringFormatWriter,
		bezeichnungStringFormatWriter: bezeichnungStringFormatWriter,
		etikettenStringFormatWriter:   etikettenStringFormatWriter,
	}
}

func (f *OrganizeNew) SetMaxKopfUndSchwanz(k, s int) {
	f.kennungStringFormatWriter.SetMaxKopfUndSchwanz(k, s)
}

func (f *OrganizeNew) WriteStringFormat(
	sw schnittstellen.WriterAndStringWriter,
	o *sku.Transacted,
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
	n2, err = f.kennungStringFormatWriter.WriteStringFormat(sw, &o.Kennung)
	n += int64(n2)
	// var n2 int64
	// n2, err = f.kennungStringFormatWriter.WriteStringFormat(
	// 	sw,
	// 	o.Kennung,
	// )
	// n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	sh := o.GetAkteSha()

	if f.options.PrintShas && (!sh.IsNull() || f.options.PrintEmptyShas) {
		n1, err = sw.WriteString("@")
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n2, err = f.shaStringFormatWriter.WriteStringFormat(sw, o.GetAkteSha())
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	t := o.GetMetadatei().GetTypPtr()

	if len(t.String()) > 0 {
		n1, err = sw.WriteString("  !")
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

	b := o.GetMetadatei().GetBezeichnungPtr()

	if f.options.PrintEtikettenAlways {
		b := o.GetMetadatei().GetEtiketten()

		for _, v := range iter.SortedValues[kennung.Etikett](b) {
			if f.options.ZittishNewlines {
				n1, err = fmt.Fprintf(sw, "\n%s", f.kennungStringFormatWriter.Padding)
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
		n1, err = fmt.Fprintf(sw, "\n%s]", f.kennungStringFormatWriter.Padding)
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

		n2, err = f.bezeichnungStringFormatWriter.WriteStringFormat(sw, b)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (f *OrganizeNew) ReadStringFormat(
	rb *catgut.RingBuffer,
	o *sku.Transacted,
) (n int64, err error) {
	if err = f.readStringFormatWithinBrackets(rb, o); err != nil {
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

func (f *OrganizeNew) readStringFormatWithinBrackets(
	rb *catgut.RingBuffer,
	o *sku.Transacted,
) (err error) {
	rr := catgut.MakeRingBufferRuneScanner(rb)

	state := 0
	var k kennung.Kennung2
	var t catgut.String
	var eof bool
	var n int

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
			if t.EqualsString("]") {
				break LOOP
			} else {
				if err = k.TodoSetBytes(&t); err != nil {
					err = errors.Wrapf(err, "Readable: %q", rb.PeekReadable())
					return
				}

				g := k.GetGattung()

				switch g {
				case gattung.Typ:
					if err = o.Metadatei.Typ.TodoSetFromKennung2(&k); err != nil {
						err = errors.Wrap(err)
						return
					}

				case gattung.Etikett:
					var e kennung.Etikett

					if err = e.TodoSetFromKennung2(&k); err != nil {
						err = errors.Wrap(err)
						return
					}

					if err = o.Metadatei.AddEtikettPtr(&e); err != nil {
						err = errors.Wrap(err)
						return
					}

				default:
					err = gattung.MakeErrUnsupportedGattung(k.GetGattung())
					return
				}

				k.Reset()
			}

		default:
			err = errors.Errorf("invalid state: %d", state)
			return
		}
	}

	// if f.options.Abbreviations.Hinweisen {
	// 	if err = f.ex.AbbreviateHinweisOnly(
	// 		&o.Kennung,
	// 	); err != nil {
	// 		err = errors.Wrap(err)
	// 		return
	// 	}
	// }

	rb.AdvanceRead(n)

	return
}
