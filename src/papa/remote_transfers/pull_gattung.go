package remote_transfers

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
)

func (c *client) PullSkus(
	ids kennung.MetaSet,
) (err error) {
	errors.TodoP1("implement etikett and akte")
	gattungInheritors := c.StoreObjekten().GetGattungInheritors(
		c,
		c,
		c.pmf,
	)

	if err = c.SkusFromFilter(
		ids,
		func(sk sku.Sku) (err error) {
			var el objekte_store.TransactedInheritor
			ok := false

			if el, ok = gattungInheritors[sk.Gattung]; !ok {
				return
			}

			errors.TodoP2("check for akte sha")

			if c.umwelt.Standort().HasObjekte(sk.Gattung, sk.ObjekteSha) {
				errors.Log().Printf("already have objekte: %s", sk.ObjekteSha)
				return
			}

			errors.Log().Printf("need objekte: %s", sk.ObjekteSha)

			if err = el.InflateFromDataIdentityAndStoreAndInherit(sk); err != nil {
				err = errors.Wrapf(err, "Sku: %s", sk)
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
