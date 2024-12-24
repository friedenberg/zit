package organize_text

import (
	"iter"
	"sort"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/tag_paths"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type (
	objSet = interfaces.MutableSetLike[*obj]
)

var objKeyer interfaces.StringKeyer[*obj]

func makeObjSet() objSet {
	return collections_value.MakeMutableValueSet(sku.GetExternalLikeKeyer[*obj]())
}

type obj struct {
	sku  sku.SkuType
	tipe tag_paths.Type
}

func (o obj) GetObjectId() *ids.ObjectId {
	return o.sku.GetObjectId()
}

func (o obj) GetSku() *sku.Transacted {
	return o.sku.GetSku()
}

func (o obj) GetSkuExternal() *sku.Transacted {
	return o.sku.GetSkuExternal()
}

func (a *obj) cloneWithType(t tag_paths.Type) (b *obj) {
	b = &obj{
		tipe: t,
		sku:  sku.CloneSkuType(a.sku),
	}

	return
}

func (a *obj) GetExternalObjectId() sku.ExternalObjectId {
	return a.sku.GetExternalObjectId()
}

func (a *obj) String() string {
	return a.sku.String()
}

func sortObjSet(
	s interfaces.MutableSetLike[*obj],
) (out []*obj) {
	out = quiter.Elements(s)

	sort.Slice(out, func(i, j int) bool {
		switch {
		case out[i].GetSkuExternal().ObjectId.IsEmpty() && out[j].GetSkuExternal().ObjectId.IsEmpty():
			return out[i].GetSkuExternal().Metadata.Description.String() < out[j].GetSkuExternal().Metadata.Description.String()

		case out[i].GetSkuExternal().ObjectId.IsEmpty():
			return true

		case out[j].GetSkuExternal().ObjectId.IsEmpty():
			return false

		default:
			return out[i].GetSkuExternal().ObjectId.String() < out[j].GetSkuExternal().ObjectId.String()
		}
	})

	return
}

type Objects []*obj

func (os Objects) Len() int {
	return len(os)
}

func (os *Objects) All() iter.Seq2[int, *obj] {
	return func(yield func(int, *obj) bool) {
		for i, o := range *os {
			if !yield(i, o) {
				break
			}
		}
	}
}

// TODO remove
func (os *Objects) Each(f interfaces.FuncIter[*obj]) (err error) {
	for _, v := range *os {
		if err = f(v); err != nil {
			if quiter.IsStopIteration(err) {
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
	sort.Slice(os, func(i, j int) bool {
		ei, ej := os[i].sku, os[j].sku

		keyI := keyer.GetKey(ei)
		keyJ := keyer.GetKey(ej)

		return keyI < keyJ
	})
}
