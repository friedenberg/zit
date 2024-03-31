package organize_text2

import (
	"fmt"
	"sort"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/erworben_cli_print_options"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

func makeObj(
	options erworben_cli_print_options.PrintOptions,
	named *sku.Transacted,
) (z *obj, err error) {
	errors.TodoP4("add bez in a better way")

	z = &obj{}

	if err = z.SetFromSkuLike(named); err != nil {
		err = errors.Wrap(err)
		return
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

	if o.Metadatei.Bezeichnung.IsEmpty() {
		return
	}

	o.Metadatei.GetEtikettenMutable().Reset()

	return
}

// TODO-P1 migrate obj to sku.Transacted
type obj struct {
	sku.Transacted
}

func (z *obj) String() string {
	return fmt.Sprintf("- [%s] %s", &z.Kennung, &z.Metadatei.Bezeichnung)
}

func sortObjSet(
	s schnittstellen.MutableSetLike[*obj],
) (out []*obj) {
	out = iter.Elements[*obj](s)

	sort.Slice(out, func(i, j int) bool {
		switch {
		case out[i].Kennung.String() != "" && out[j].Kennung.String() != "":
			return out[i].Kennung.String() < out[j].Kennung.String()

		case out[i].Kennung.String() == "":
			return true

		case out[j].Kennung.String() == "":
			return false

		default:
			return out[i].Metadatei.Bezeichnung.String() < out[j].Metadatei.Bezeichnung.String()
		}
	})

	return
}
