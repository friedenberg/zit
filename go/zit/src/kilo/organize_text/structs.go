package organize_text

import (
	"fmt"
	"sort"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

// TODO-P1 migrate obj to sku.Transacted
type obj struct {
	sku.Transacted
	virtual bool
}

func (z *obj) String() string {
	return fmt.Sprintf("- [%s] %s", &z.Kennung, &z.Metadatei.Bezeichnung)
}

func sortObjSet(
	s schnittstellen.MutableSetLike[*obj],
) (out []*obj) {
	out = iter.Elements(s)

	sort.Slice(out, func(i, j int) bool {
		switch {
		case out[i].Kennung.IsEmpty() && out[j].Kennung.IsEmpty():
			return out[i].Metadatei.Bezeichnung.String() < out[j].Metadatei.Bezeichnung.String()

		case out[i].Kennung.IsEmpty():
			return true

		case out[j].Kennung.IsEmpty():
			return false

		default:
			return out[i].Kennung.String() < out[j].Kennung.String()
		}
	})

	return
}

type Objekten []*obj

func (os Objekten) Len() int {
	return len(os)
}

func (os *Objekten) Each(f schnittstellen.FuncIter[*obj]) (err error) {
	for _, v := range *os {
		if err = f(v); err != nil {
			if iter.IsStopIteration(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (os Objekten) Any() *obj {
	for _, v := range os {
		return v
	}

	return nil
}

func (os *Objekten) Add(v *obj) error {
	*os = append(*os, v)
	return nil
}

func (os *Objekten) Del(v *obj) error {
	for i, v1 := range *os {
		if v == v1 {
			*os = append((*os)[:i], (*os)[i+1:]...)
			break
		}
	}

	return nil
}

func (os Objekten) Sort() {
	out := os

	sort.Slice(out, func(i, j int) bool {
		switch {
		case out[i].Kennung.IsEmpty() && out[j].Kennung.IsEmpty():
			return out[i].Metadatei.Bezeichnung.String() < out[j].Metadatei.Bezeichnung.String()

		case out[i].Kennung.IsEmpty():
			return true

		case out[j].Kennung.IsEmpty():
			return false

		default:
			return out[i].Kennung.String() < out[j].Kennung.String()
		}
	})
}
