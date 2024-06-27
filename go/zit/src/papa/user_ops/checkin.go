package user_ops

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

type Checkin struct {
	Delete bool
}

func (c Checkin) Run(
	u *umwelt.Umwelt,
	qg query.GroupWithKasten,
) (err error) {
	u.Lock()
	defer errors.Deferred(&err, u.Unlock)

	ui.Log().Print(qg)

	if err = u.GetStore().ReadExternal(
		qg,
		func(col sku.CheckedOutLike) (err error) {
			if _, err = u.GetStore().CreateOrUpdateCheckedOut(
				col,
				true,
			); err != nil {
				ui.Debug().Print(err)
				err = errors.Wrap(err)
				return
			}

			if !c.Delete {
				return
			}

			if err = u.GetStore().DeleteCheckout(col); err != nil {
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
