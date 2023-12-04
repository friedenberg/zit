package sku_fmt

import (
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/erworben_cli_print_options"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/catgut"
	"github.com/friedenberg/zit/src/delta/thyme"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type organizeNew struct {
	options erworben_cli_print_options.PrintOptions

	shaStringFormatWriter         schnittstellen.StringFormatWriter[schnittstellen.ShaLike]
	kennungStringFormatWriter     schnittstellen.StringFormatWriter[*kennung.Kennung2]
	typStringFormatWriter         schnittstellen.StringFormatWriter[*kennung.Typ]
	bezeichnungStringFormatWriter schnittstellen.StringFormatWriter[*bezeichnung.Bezeichnung]
	etikettenStringFormatWriter   schnittstellen.StringFormatWriter[kennung.EtikettSet]

	ex kennung.Abbr

	maxKopf, maxSchwanz int
}

func MakeOrganizeNewFormat(
	options erworben_cli_print_options.PrintOptions,
	shaStringFormatWriter schnittstellen.StringFormatWriter[schnittstellen.ShaLike],
	kennungStringFormatWriter schnittstellen.StringFormatWriter[*kennung.Kennung2],
	typStringFormatWriter schnittstellen.StringFormatWriter[*kennung.Typ],
	bezeichnungStringFormatWriter schnittstellen.StringFormatWriter[*bezeichnung.Bezeichnung],
	etikettenStringFormatWriter schnittstellen.StringFormatWriter[kennung.EtikettSet],
) *organizeNew {
	options.PrintTime = false
	options.PrintShas = false

	return &organizeNew{
		options:                       options,
		shaStringFormatWriter:         shaStringFormatWriter,
		kennungStringFormatWriter:     kennungStringFormatWriter,
		typStringFormatWriter:         typStringFormatWriter,
		bezeichnungStringFormatWriter: bezeichnungStringFormatWriter,
		etikettenStringFormatWriter:   etikettenStringFormatWriter,
	}
}

func (f *organizeNew) SetMaxKopf(m int) {
	f.maxKopf = m
}

func (f *organizeNew) SetMaxSchwanz(m int) {
	f.maxSchwanz = m
}

func (f *organizeNew) WriteStringFormat(
	sw schnittstellen.WriterAndStringWriter,
	o *sku.Transacted,
) (n int64, err error) {
	var n1 int

	if f.options.PrintTime {
		t := o.GetTai()

		n1, err = sw.WriteString(t.Format(thyme.FormatDateTime))
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

	h := kennung.Aligned(&o.Kennung, f.maxKopf, f.maxSchwanz)
	n1, err = sw.WriteString(h)
	n += int64(n1)
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

	var n2 int64

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
		n1, err = sw.WriteString(" !")
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

		if b.Len() > 0 {
			n1, err = sw.WriteString(" ")
			n += int64(n1)

			if err != nil {
				err = errors.Wrap(err)
				return
			}

			n2, err = f.etikettenStringFormatWriter.WriteStringFormat(sw, b)
			n += n2

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	n1, err = sw.WriteString("]")
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

func (f *organizeNew) ReadStringFormat(
	rb io.Reader,
	o *sku.Transacted,
) (n int64, err error) {
	cg := catgut.GetPool().Get()

	if _, err = io.Copy(cg, rb); err != nil {
		err = errors.Wrap(err)
		return
	}

	remaining := cg.String()

	if len(remaining) < 3 {
		err = errors.Errorf("expected at least 3 characters")
		return
	}

	if remaining[:3] != "- [" {
		err = errors.Errorf("expected '- [', but got '%s'", remaining[:3])
		return
	}

	remaining = remaining[3:]

	idx := -1

	if idx = strings.Index(remaining, "]"); idx == -1 {
		err = errors.Errorf("expected ']' after hinweis, but not found")
		return
	}

	remainingKennung := strings.TrimSpace(remaining[:idx])

	if idxSpace := strings.Index(remainingKennung, " "); idxSpace != -1 {
		remainingKennung = remainingKennung[:idxSpace]
	}

	if err = o.Kennung.Set(remainingKennung); err != nil {
		err = errors.Wrap(err)
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

	// no bezeichnung
	if idx+2 > len(remaining)-1 {
		return
	}

	remaining = remaining[idx+2:]

	if err = o.Metadatei.Bezeichnung.Set(remaining); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
