package sku_fmt

import (
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/erworben_cli_print_options"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/catgut"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type KennungAlignedFormat interface {
	SetMaxKopf(m int)
	SetMaxSchwanz(m int)
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

func (f *organize) SetMaxKopf(m int) {
	f.maxKopf = m
}

func (f *organize) SetMaxSchwanz(m int) {
	f.maxSchwanz = m
}

func (f *organize) ReadStringFormat(
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
