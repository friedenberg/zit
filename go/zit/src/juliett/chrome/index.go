package chrome

import (
	"net/url"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

// TODO abstract and regenerate on commit / reindex
func (c *Store) initializeIndex() (err error) {
	if err = c.initializeCache(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var l sync.Mutex

	if err = c.externalStoreInfo.FuncPrimitiveQuery(
		nil,
		func(sk *sku.Transacted) (err error) {
			if !sk.GetType().Equals(c.typ) {
				return
			}

			var u *url.URL

			if u, err = c.getUrl(sk); err != nil {
				err = errors.Wrap(err)
				return
			}

			cl := sku.GetTransactedPool().Get()
			sku.TransactedResetter.ResetWith(cl, sk)

			l.Lock()
			defer l.Unlock()

			{
				existing, ok := c.transactedUrlIndex[*u]

				if !ok {
					existing = sku.MakeTransactedMutableSet()
					c.transactedUrlIndex[*u] = existing
				}

				if err = existing.Add(cl); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			{
				existing, ok := c.tabCache.Rows[sk.ObjectId.String()]

				if ok {
					c.transactedTabIdIndex[existing] = cl
				}
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
