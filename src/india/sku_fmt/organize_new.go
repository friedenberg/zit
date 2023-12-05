package sku_fmt

import (
	"bufio"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/erworben_cli_print_options"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/zittish"
	"github.com/friedenberg/zit/src/charlie/catgut"
	"github.com/friedenberg/zit/src/charlie/gattung"
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
	rb *catgut.RingBuffer,
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

	if tokens, err = f.readStringFormatWithinBrackets(tokens, o); err != nil {
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

func (f *organizeNew) readStringFormatWithinBrackets(
	tokens []string,
	o *sku.Transacted,
) (remainingTokens []string, err error) {
	remainingTokens = tokens

	state := 0
	var k kennung.Kennung2
	var i int
	var t string

LOOP:
	for i, t = range tokens {
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
			if t == "]" {
				break LOOP
			} else {
				if err = k.Set(t); err != nil {
					err = errors.Wrap(err)
					return
				}

				g := k.GetGattung()

				switch g {
				case gattung.Typ:
					if err = o.Metadatei.Typ.Set(k.String()); err != nil {
						err = errors.Wrap(err)
						return
					}

				case gattung.Etikett:
					var e kennung.Etikett

					if err = e.Set(k.String()); err != nil {
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

	return
}
