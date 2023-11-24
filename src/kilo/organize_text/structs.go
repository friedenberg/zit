package organize_text

import (
	"fmt"
	"sort"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/erworben_cli_print_options"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

func makeObj(
	options erworben_cli_print_options.PrintOptions,
	named *sku.Transacted,
	expanders kennung.Abbr,
) (z obj, err error) {
	errors.TodoP4("add bez in a better way")
	z = obj{
		Sku: *named,
	}

	if options.Abbreviations.Hinweisen {
		if z.Sku.Kennung, err = expanders.AbbreviateHinweisOnly(
			z.Sku.Kennung,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

type obj struct {
	Sku   sku.Transacted
	IsNew bool
}

func (a obj) Len() int {
	return len(a.Sku.Kennung.String())
}

func (a obj) LenKopfUndSchwanz() (int, int) {
	kopf, schwanz := a.KopfUndSchwanz()

	return len(kopf), len(schwanz)
}

func (a obj) KopfUndSchwanz() (kopf, schwanz string) {
	parts := a.Sku.Kennung.Parts()
	kopf = parts[0]
	schwanz = parts[2]

	return
}

func (a obj) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a obj) Equals(b obj) bool {
	if a.Sku.Kennung.String() != b.Sku.Kennung.String() {
		return false
	}

	if !a.Sku.Metadatei.Bezeichnung.Equals(b.Sku.Metadatei.Bezeichnung) {
		return false
	}

	return true
}

func (z obj) String() string {
	return fmt.Sprintf("- [%s] %s", z.Sku.Kennung, z.Sku.Metadatei.Bezeichnung)
}

func (z *obj) setExistingObj(
	options erworben_cli_print_options.PrintOptions,
	v string,
	ex kennung.Abbr,
) (err error) {
	remaining := v

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

	if err = z.Sku.Kennung.Set(remainingKennung); err != nil {
		err = errors.Wrap(err)
		return
	}

	if options.Abbreviations.Hinweisen {
		if z.Sku.Kennung, err = ex.AbbreviateHinweisOnly(
			z.Sku.Kennung,
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

	if err = z.Sku.Metadatei.Bezeichnung.Set(remaining); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (z *obj) setNewObj(v string) (err error) {
	remaining := v

	if remaining[:2] != "- " {
		err = errors.Errorf("expected '- ', but got '%s'", remaining[:2])
		return
	}

	remaining = remaining[2:]

	if err = z.Sku.Metadatei.Bezeichnung.Set(remaining); err != nil {
		err = errors.Wrap(err)
		return
	}

	z.IsNew = true

	return
}

func sortObjSet(
	s schnittstellen.MutableSetLike[obj],
) (out []obj) {
	out = iter.Elements[obj](s)

	sort.Slice(out, func(i, j int) bool {
		switch {
		case out[i].Sku.Kennung.String() != "" && out[j].Sku.Kennung.String() != "":
			return out[i].Sku.Kennung.String() < out[j].Sku.Kennung.String()

		case out[i].Sku.Kennung.String() == "":
			return true

		case out[j].Sku.Kennung.String() == "":
			return false

		default:
			return out[i].Sku.Metadatei.Bezeichnung.String() < out[j].Sku.Metadatei.Bezeichnung.String()
		}
	})

	return
}
