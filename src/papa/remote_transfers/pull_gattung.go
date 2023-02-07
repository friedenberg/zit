package remote_transfers

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
)

func (c *client) PullSkus(
	filter kennung.Filter,
	gattungSet gattungen.Set,
) (err error) {
	errors.TodoP0("implement etikett and akte")
	gattungInheritors := map[gattung.Gattung]objekte_store.TransactedInheritor{
		gattung.Zettel: c.GetInheritorZettel(c, c),
		gattung.Typ:    c.GetInheritorTyp(c, c),
	}

	if err = c.SkusFromFilter(
		filter,
		gattungSet,
		func(sk sku.Sku2) (err error) {
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
