package organize_text

import (
	"sort"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/tag_paths"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type (
	objSet = interfaces.MutableSetLike[*obj]
)

var objKeyer interfaces.StringKeyer[*obj]

func makeObjSet() objSet {
	return collections_value.MakeMutableValueSet(objKeyer)
}

// TODO-P1 migrate obj to sku.Transacted
type obj struct {
	sku.Transacted
	tag_paths.Type
}

func (a *obj) cloneWithType(t tag_paths.Type) (b *obj) {
	b = &obj{Type: t}
	sku.TransactedResetter.ResetWith(&b.Transacted, &a.Transacted)
	return
}

func (a *obj) String() string {
	return a.Transacted.String()
}

func sortObjSet(
	s interfaces.MutableSetLike[*obj],
) (out []*obj) {
	out = iter.Elements(s)

	sort.Slice(out, func(i, j int) bool {
		switch {
		case out[i].ObjectId.IsEmpty() && out[j].ObjectId.IsEmpty():
			return out[i].Metadata.Description.String() < out[j].Metadata.Description.String()

		case out[i].ObjectId.IsEmpty():
			return true

		case out[j].ObjectId.IsEmpty():
			return false

		default:
			return out[i].ObjectId.String() < out[j].ObjectId.String()
		}
	})

	return
}

type Objekten []*obj

func (os Objekten) Len() int {
	return len(os)
}

func (os *Objekten) Each(f interfaces.FuncIter[*obj]) (err error) {
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
		case out[i].ObjectId.IsEmpty() && out[j].ObjectId.IsEmpty():
			return out[i].Metadata.Description.String() < out[j].Metadata.Description.String()

		case out[i].ObjectId.IsEmpty():
			return true

		case out[j].ObjectId.IsEmpty():
			return false

		default:
			// TODO sort by ints for virtual kennung
			return out[i].ObjectId.String() < out[j].ObjectId.String()
		}
	})
}
