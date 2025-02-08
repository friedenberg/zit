package organize_text

import (
	"fmt"
	"iter"
	"sort"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/options_print"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
)

func MakeSkuMapWithOrder(c int) (out SkuMapWithOrder) {
	out.m = make(map[string]skuTypeWithIndex, c)
	return
}

type skuTypeWithIndex struct {
	sku sku.SkuType
	int
}

type SkuMapWithOrder struct {
	m    map[string]skuTypeWithIndex
	next int
}

func (smwo *SkuMapWithOrder) AllSkuAndIndex() iter.Seq2[int, sku.SkuType] {
	return func(yield func(int, sku.SkuType) bool) {
		for _, sk := range smwo.m {
			if !yield(sk.int, sk.sku) {
				break
			}
		}
	}
}

func (smwo *SkuMapWithOrder) AsExternalLikeSet() sku.SkuTypeSetMutable {
	elms := sku.MakeSkuTypeSetMutable()
	errors.PanicIfError(smwo.Each(elms.Add))
	return elms
}

func (smwo *SkuMapWithOrder) AsTransactedSet() sku.TransactedMutableSet {
	tms := sku.MakeTransactedMutableSet()
	errors.PanicIfError(smwo.Each(func(el sku.SkuType) (err error) {
		return tms.Add(el.GetSkuExternal())
	}))
	return tms
}

func (sm *SkuMapWithOrder) Del(sk sku.SkuType) error {
	delete(sm.m, keyer.GetKey(sk))
	return nil
}

func (sm *SkuMapWithOrder) Add(sk sku.SkuType) error {
	k := keyer.GetKey(sk)
	entry, ok := sm.m[k]

	if !ok {
		entry.int = sm.next
		entry.sku = sk
		sm.next++
	}

	sm.m[k] = entry

	return nil
}

func (sm *SkuMapWithOrder) Len() int {
	return len(sm.m)
}

func (sm *SkuMapWithOrder) Clone() (out SkuMapWithOrder) {
	out = MakeSkuMapWithOrder(sm.Len())

	for _, v := range sm.m {
		out.Add(v.sku)
	}

	return out
}

func (sm SkuMapWithOrder) Sorted() (out []sku.SkuType) {
	out = make([]sku.SkuType, 0, sm.Len())

	for _, v := range sm.m {
		out = append(out, v.sku)
	}

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

func (sm *SkuMapWithOrder) Each(
	f interfaces.FuncIter[sku.SkuType],
) (err error) {
	for _, v := range sm.Sorted() {
		_, ok := sm.m[keyer.GetKey(v)]

		if !ok {
			continue
		}

		if err = f(v); err != nil {
			if errors.IsStopIteration(err) {
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
	Before, After  SkuMapWithOrder
	Added, Removed SkuMapWithOrder
	Changed        SkuMapWithOrder
}

func (c Changes) String() string {
	return fmt.Sprintf(
		"Before: %d, After: %d, Added: %d, Removed: %d, Changed: %d",
		c.Before.Len(),
		c.After.Len(),
		c.Added.Len(),
		c.Removed.Len(),
		c.Changed.Len(),
	)
}

// TODO combine with above
type OrganizeResults struct {
	Before, After *Text
	Original      sku.SkuTypeSet
	QueryGroup    *query.Query
}

func ChangesFrom(
	po options_print.V0,
	a, b *Text,
	original sku.SkuTypeSet,
) (c Changes, err error) {
	if c, err = ChangesFromResults(
		po,
		OrganizeResults{
			Before:   a,
			After:    b,
			Original: original,
		}); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func ChangesFromResults(
	po options_print.V0,
	results OrganizeResults,
) (c Changes, err error) {
	if err = applyToText(po, results.Before); err != nil {
		err = errors.Wrap(err)
		return
	}

	if c.Before, err = results.Before.GetSkus(results.Original); err != nil {
		err = errors.Wrap(err)
		return
	}

	if c.After, err = results.After.GetSkus(results.Original); err != nil {
		err = errors.Wrap(err)
		return
	}

	c.Changed = c.After.Clone()
	c.Removed = c.Before.Clone()

	for _, sk := range c.After.m {
		if err = c.Removed.Del(sk.sku); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	for _, sk := range c.Removed.AllSkuAndIndex() {
		if err = results.Before.RemoveFromTransacted(sk); err != nil {
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

func applyToText(
	po options_print.V0,
	t *Text,
) (err error) {
	if po.PrintTagsAlways {
		return
	}

	if err = t.Options.Skus.Each(
		func(el sku.SkuType) (err error) {
			sk := el.GetSkuExternal()

			if sk.Metadata.Description.IsEmpty() {
				return
			}

			sk.Metadata.ResetTags()

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
