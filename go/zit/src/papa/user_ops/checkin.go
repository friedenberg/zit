package user_ops

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type Checkin struct {
	Delete bool
}

func (c Checkin) Run(
	u *env.Env,
	eqwk *query.Group,
) (err error) {
	u.Lock()
	defer errors.Deferred(&err, u.Unlock)

	ui.Log().Print(eqwk)

	if err = u.GetStore().QueryCheckedOut(
		eqwk,
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
