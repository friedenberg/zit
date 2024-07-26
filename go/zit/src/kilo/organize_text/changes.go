package organize_text

import (
	"sort"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func MakeSkuMapWithOrder(c int) (out SkuMapWithOrder) {
	out.m = make(map[string]sku.ExternalLike, c)
	out.o = make([]sku.ExternalLike, 0, c)
	return
}

type SkuMapWithOrder struct {
	m map[string]sku.ExternalLike
	o []sku.ExternalLike
}

func (sm *SkuMapWithOrder) Del(sk sku.ExternalLike) error {
	delete(sm.m, key(sk))
	return nil
}

func (sm *SkuMapWithOrder) Add(sk sku.ExternalLike) error {
	k := key(sk)
	_, ok := sm.m[k]

	if !ok {
		sm.m[k] = sk
		sm.o = append(sm.o, sk)
	}

	return nil
}

func (sm *SkuMapWithOrder) Len() int {
	return len(sm.m)
}

func (sm *SkuMapWithOrder) Clone() (out SkuMapWithOrder) {
	out = MakeSkuMapWithOrder(sm.Len())

	for _, v := range sm.m {
		out.Add(v)
	}

	return out
}

func (sm SkuMapWithOrder) Sort() {
	out := sm.o

	sort.Slice(out, func(i, j int) bool {
		switch {
		case out[i].GetSku().ObjectId.IsEmpty() && out[j].GetSku().ObjectId.IsEmpty():
			return out[i].GetSku().Metadata.Description.String() < out[j].GetSku().Metadata.Description.String()

		case out[i].GetSku().ObjectId.IsEmpty():
			return true

		case out[j].GetSku().ObjectId.IsEmpty():
			return false

		default:
			return out[i].GetSku().ObjectId.String() < out[j].GetSku().ObjectId.String()
		}
	})
}

func (sm *SkuMapWithOrder) Each(
	f interfaces.FuncIter[sku.ExternalLike],
) (err error) {
	sm.Sort()

	for _, v := range sm.o {
		_, ok := sm.m[key(v)]

		if !ok {
			continue
		}

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

type Changes struct {
	a, b           SkuMapWithOrder
	added, removed SkuMapWithOrder
	Changed        SkuMapWithOrder
}

func ChangesFrom(
	a, b *Text,
	original sku.ExternalLikeSet,
) (c Changes, err error) {
	if c.a, err = a.GetSkus(original); err != nil {
		err = errors.Wrap(err)
		return
	}

	if c.b, err = b.GetSkus(original); err != nil {
		err = errors.Wrap(err)
		return
	}

	c.Changed = c.b.Clone()
	c.removed = c.a.Clone()

	for _, sk := range c.b.o {
		if err = c.removed.Del(sk); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	for _, sk := range c.removed.o {
		if err = a.RemoveFromTransacted(sk); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = c.Changed.Add(sk); err != nil {
			err = errors.Wrap(err)
			return
		}

	}

	return
}
