package organize_text

import (
	"fmt"
	"sort"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/erworben_cli_print_options"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

func makeObj(
	options erworben_cli_print_options.PrintOptions,
	named *sku.Transacted,
	expanders kennung.Abbr,
) (z *obj, err error) {
	errors.TodoP4("add bez in a better way")

	z = &obj{}

	if err = z.Sku.SetFromSkuLike(named); err != nil {
		err = errors.Wrap(err)
		return
	}

	if options.Abbreviations.Hinweisen {
		if err = expanders.AbbreviateHinweisOnly(
			&z.Sku.Kennung,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = z.removeEtikettenIfNecessary(options); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (o *obj) removeEtikettenIfNecessary(
	options erworben_cli_print_options.PrintOptions,
) (err error) {
	if options.PrintEtikettenAlways {
		return
	}

	if o.Sku.Metadatei.Bezeichnung.IsEmpty() {
		return
	}

	o.Sku.Metadatei.GetEtikettenMutable().Reset()

	return
}

// TODO-P1 migrate obj to sku.Transacted
type obj struct {
	Sku sku.Transacted
}

func (z *obj) String() string {
	return fmt.Sprintf("- [%s] %s", &z.Sku.Kennung, &z.Sku.Metadatei.Bezeichnung)
}

func sortObjSet(
	s schnittstellen.MutableSetLike[*obj],
) (out []*obj) {
	out = iter.Elements[*obj](s)

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
