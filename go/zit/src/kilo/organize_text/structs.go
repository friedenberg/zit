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
	objSet  = interfaces.MutableSetLike[*obj]
	skuType = sku.ExternalLike
)

var objKeyer interfaces.StringKeyer[*obj]

func makeObjSet() objSet {
	return collections_value.MakeMutableValueSet(objKeyer)
}

// TODO-P1 migrate obj to sku.Transacted
type obj struct {
	Transacted skuType
	tag_paths.Type
}

func (a *obj) cloneWithType(t tag_paths.Type) (b *obj) {
	b = &obj{Type: t, Transacted: a.Transacted.Clone()}
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
		case out[i].Transacted.GetSku().ObjectId.IsEmpty() && out[j].Transacted.GetSku().ObjectId.IsEmpty():
			return out[i].Transacted.GetSku().Metadata.Description.String() < out[j].Transacted.GetSku().Metadata.Description.String()

		case out[i].Transacted.GetSku().ObjectId.IsEmpty():
			return true

		case out[j].Transacted.GetSku().ObjectId.IsEmpty():
			return false

		default:
			return out[i].Transacted.GetSku().ObjectId.String() < out[j].Transacted.GetSku().ObjectId.String()
		}
	})

	return
}

type Objects []*obj

func (os Objects) Len() int {
	return len(os)
}

func (os *Objects) Each(f interfaces.FuncIter[*obj]) (err error) {
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

func (os Objects) Any() *obj {
	for _, v := range os {
		return v
	}

	return nil
}

func (os *Objects) Add(v *obj) error {
	*os = append(*os, v)
	return nil
}

func (os *Objects) Del(v *obj) error {
	for i, v1 := range *os {
		if v == v1 {
			*os = append((*os)[:i], (*os)[i+1:]...)
			break
		}
	}

	return nil
}

func (os Objects) Sort() {
	out := os

	sort.Slice(out, func(i, j int) bool {
		switch {
		case out[i].Transacted.GetSku().ObjectId.IsEmpty() && out[j].Transacted.GetSku().ObjectId.IsEmpty():
			return out[i].Transacted.GetSku().Metadata.Description.String() < out[j].Transacted.GetSku().Metadata.Description.String()

		case out[i].Transacted.GetSku().ObjectId.IsEmpty():
			return true

		case out[j].Transacted.GetSku().ObjectId.IsEmpty():
			return false

		default:
			// TODO sort by ints for virtual object id
			return out[i].Transacted.GetSku().ObjectId.String() < out[j].Transacted.GetSku().ObjectId.String()
		}
	})
}
