package changes2

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type Changes2 struct {
	A, B           sku.TransactedSet
	Added, Removed sku.TransactedMutableSet
}

func ChangesFrom2(
	a, b Changeable,
	original sku.TransactedSet,
) (c Changes2, err error) {
	if c.A, err = a.GetSkus(original); err != nil {
		err = errors.Wrap(err)
		return
	}

	if c.B, err = b.GetSkus(original); err != nil {
		err = errors.Wrap(err)
		return
	}

	c.Added = sku.MakeTransactedMutableSet()
	c.Removed = c.A.CloneMutableSetLike()

	if err = c.B.Each(
		func(sk *sku.Transacted) (err error) {
			if !c.A.Contains(sk) {
				if err = c.Added.Add(sk); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			if err = c.Removed.Del(sk); err != nil {
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
