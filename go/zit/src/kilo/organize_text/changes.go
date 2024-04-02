package organize_text

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type Changes struct {
	a, b           sku.TransactedSet
	added, removed sku.TransactedMutableSet
	Changed        sku.TransactedMutableSet
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

	c.Changed = c.b.CloneMutableSetLike()
	c.removed = c.a.CloneMutableSetLike()

	if err = c.b.Each(
		func(sk *sku.Transacted) (err error) {
			if err = c.removed.Del(sk); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.removed.Each(
		func(sk *sku.Transacted) (err error) {
			if err = a.RemoveFromTransacted(sk); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = c.Changed.Add(sk); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
