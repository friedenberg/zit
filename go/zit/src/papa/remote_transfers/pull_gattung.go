package remote_transfers

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/juliett/query"
)

func (c *client) PullSkus(
	ids *query.Group,
) (err error) {
	errors.TodoP1("implement etikett and akte")
	// if err = c.SkusFromFilter(
	// 	ids,
	// 	func(sk *sku.Transacted) (err error) {
	// 		var el objekte_store.TransactedInheritor
	// 		ok := false

	// 		if el, ok = gattungInheritors[gattung.Must(sk.GetGattung())]; !ok {
	// 			return
	// 		}

	// 		errors.TodoP2("check for akte sha")

	// 		if c.umwelt.Standort().HasObjekte(
	// 			c.umwelt.Konfig().GetStoreVersion(),
	// 			sk.GetGattung(),
	// 			sk.GetObjekteSha(),
	// 		) {
	// 			errors.Log().Printf("already have objekte: %s", sk.GetObjekteSha())
	// 			return
	// 		}

	// 		if err = el.InflateFromDataIdentityAndStoreAndInherit(sk); err != nil {
	// 			err = errors.Wrapf(err, "Sku: %s", sk)
	// 			return
	// 		}

	// 		return
	// 	},
	// ); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	return
}
