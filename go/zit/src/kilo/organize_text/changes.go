package organize_text

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type SkuMap map[string]*sku.Transacted

func (sm *SkuMap) Del(sk *sku.Transacted) error {
	delete((*sm), key(sk))
	return nil
}

func (sm *SkuMap) Add(sk *sku.Transacted) error {
	(*sm)[key(sk)] = sk
	return nil
}

func (sm *SkuMap) Len() int {
	return len(*sm)
}

func (sm *SkuMap) Clone() SkuMap {
	out := make(SkuMap, sm.Len())

	for k, v := range *sm {
		out[k] = v
	}

	return out
}

type Changes struct {
	a, b           SkuMap
	added, removed SkuMap
	Changed        SkuMap
}

func ChangesFrom(
	a, b *Text,
	original sku.TransactedSet,
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

	for _, sk := range c.b {
		if err = c.removed.Del(sk); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	for _, sk := range c.removed {
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
