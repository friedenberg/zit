package sku_fmt

import (
	"bufio"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/erworben_cli_print_options"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/zittish"
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
	rb io.Reader,
	o *sku.Transacted,
) (n int64, err error) {
	scanner := bufio.NewScanner(rb)

	scanner.Split(zittish.SplitMatcher)

	tokens := make([]string, 0)

	beforeHyphen := false

	for scanner.Scan() {
		t := scanner.Text()

		if t == " " && beforeHyphen {
			continue
		}

		tokens = append(tokens, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if len(tokens) < 1 {
		err = errors.Errorf("no tokens")
		return
	}

	if tokens[0] != "-" {
		err = errors.Errorf("expected %q at beginning but to got %q", "-", tokens[0])
		return
	}

	tokens = tokens[1:]

	if tokens, err = f.readStringFormatWithKennung(tokens, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	remaining := strings.Join(tokens, "")

	if err = o.Metadatei.Bezeichnung.Set(remaining); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *organize) readStringFormatWithKennung(
	tokens []string,
	o *sku.Transacted,
) (remainingTokens []string, err error) {
	remainingTokens = tokens

	state := 0

LOOP:
	for i, t := range tokens {
		if t == " " {
			continue
		}

		switch state {
		case 0:
			if t != "[" {
				return
			}

			state++

		case 1:
			if err = o.Kennung.Set(t); err != nil {
				o.Kennung.Reset()
				return
			}

			state++

		case 2:
			if t != "]" {
				o.Kennung.Reset()
				return
			}

			if len(remainingTokens) > i {
				remainingTokens = tokens[i+1:]
			} else {
				remainingTokens = nil
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

	return
}
