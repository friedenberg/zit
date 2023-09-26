package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type Merge struct {
	*umwelt.Umwelt
}

func (op Merge) Run(skus sku.TransactedHeap) (err error) {
	// TODO-P1 make data structure for merging
	q := op.MakeQueryAll()

	if err = q.Set(":"); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = op.StoreObjekten().Query(
		q,
		func(sk *sku.Transacted) (err error) {
			peeked, ok := skus.Peek()

			if !ok {
				return iter.MakeErrStopIteration()
			}

			if peeked.EqualsSkuLikePtr(sk) {
				skus.Pop()
				return
			}

			if sk.GetTai().Less(peeked.GetTai()) {
				return
			}

			_, _ = skus.PopAndSave()

			// return errors.Normalf(
			// 	"merge required: %q < %q",
			// 	sku_formats.String(peeked),
			// 	sku_formats.String(sk),
			// )

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	for {
		_, ok := skus.PopAndSave()

		if !ok {
			break
		}
	}

	skus.Restore()

	return
}
