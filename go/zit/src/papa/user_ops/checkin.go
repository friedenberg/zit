package user_ops

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type Checkin struct {
	Delete   bool
	Organize bool
}

func (c Checkin) Run(
	u *env.Env,
	qg *query.Group,
) (err error) {
	if c.Organize {
		if err = c.runOrganize(u, qg); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = c.runNoOrganize(u, qg); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (c Checkin) runOrganize(
	u *env.Env,
	qg *query.Group,
) (err error) {
	opOrganize := Organize{
		Env: u,
	}

	ui.Log().Print(qg)

	if err = opOrganize.RunWithQueryGroup(
		qg,
		func(changed sku.ExternalLike) (err error) {
			if err = u.GetStore().CreateOrUpdate(
				changed,
				objekte_mode.Make(
					objekte_mode.ModeMergeCheckedOut,
				),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

      // TODO mark for deletion

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Checkin) runNoOrganize(
	u *env.Env,
	qg *query.Group,
) (err error) {
	u.Lock()
	defer errors.Deferred(&err, u.Unlock)

	ui.Log().Print(qg)

	if err = u.GetStore().QueryCheckedOut(
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
