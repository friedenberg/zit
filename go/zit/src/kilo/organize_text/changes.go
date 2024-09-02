package organize_text

import (
	"sort"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/erworben_cli_print_options"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
)

func MakeSkuMapWithOrder(c int) (out SkuMapWithOrder) {
	out.m = make(map[string]skuExternalLikeWithIndex, c)
	return
}

type skuExternalLikeWithIndex struct {
	sku.ExternalLike
	int
}

type SkuMapWithOrder struct {
	m    map[string]skuExternalLikeWithIndex
	next int
}

func (smwo *SkuMapWithOrder) AsExternalLikeSet() sku.ExternalLikeMutableSet {
	elms := sku.MakeExternalLikeMutableSet()
	errors.PanicIfError(smwo.Each(elms.Add))
	return elms
}

func (sm *SkuMapWithOrder) Del(sk sku.ExternalLike) error {
	delete(sm.m, key(sk))
	return nil
}

func (sm *SkuMapWithOrder) Add(sk sku.ExternalLike) error {
	k := key(sk)
	entry, ok := sm.m[k]

	if !ok {
		entry.int = sm.next
		entry.ExternalLike = sk
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
		out.Add(v)
	}

	return out
}

func (sm SkuMapWithOrder) Sorted() (out []sku.ExternalLike) {
	out = make([]sku.ExternalLike, 0, sm.Len())

	for _, v := range sm.m {
		out = append(out, v.ExternalLike)
	}

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

	return
}

func (sm *SkuMapWithOrder) Each(
	f interfaces.FuncIter[sku.ExternalLike],
) (err error) {
	for _, v := range sm.Sorted() {
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
	Before, After  SkuMapWithOrder
	Added, Removed SkuMapWithOrder
	Changed        SkuMapWithOrder
}

// TODO combine with above
type OrganizeResults struct {
	Before, After *Text
	Original      sku.ExternalLikeSet
	QueryGroup    *query.Group
}

func ChangesFrom(
	po erworben_cli_print_options.PrintOptions,
	a, b *Text,
	original sku.ExternalLikeSet,
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
	po erworben_cli_print_options.PrintOptions,
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
		if err = c.Removed.Del(sk); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	for _, sk := range c.Removed.m {
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
	po erworben_cli_print_options.PrintOptions,
	t *Text,
) (err error) {
	if po.PrintTagsAlways {
		return
	}

	if err = t.Options.Skus.Each(
		func(el sku.ExternalLike) (err error) {
			sk := el.GetSku()

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
