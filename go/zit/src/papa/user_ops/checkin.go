package user_ops

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/organize_text"
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

	var organizeResults organize_text.OrganizeResults

	if organizeResults, err = opOrganize.RunWithQueryGroup(
		qg,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	u.Lock()
	defer errors.Deferred(&err, u.Unlock)

	if _, err = u.CommitRemainingOrganizeResults(organizeResults); err != nil {
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
			if err = u.GetStore().CreateOrUpdateCheckedOut(
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
