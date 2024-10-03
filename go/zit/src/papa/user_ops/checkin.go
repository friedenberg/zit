package user_ops

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
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

func (op Checkin) Run(
	u *env.Env,
	qg *query.Group,
) (err error) {
	if op.Organize {
		if err = op.runOrganize(u, qg); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = u.Lock(); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = u.GetStore().QueryCheckedOut(
			qg,
			func(col sku.CheckedOutLike) (err error) {
				if err = u.GetStore().CreateOrUpdateCheckedOut(
					col,
					!op.Delete,
				); err != nil {
					err = errors.Wrap(err)
					return
				}

				if !op.Delete {
					return
				}

				if err = u.GetStore().DeleteCheckedOutLike(col); err != nil {
					err = errors.Wrap(err)
					return
				}

				return
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = u.Unlock(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (op Checkin) runOrganize(
	u *env.Env,
	qg *query.Group,
) (err error) {
	flagDelete := organize_text.OptionCommentBooleanFlag{
		Value:   &op.Delete,
		Comment: "delete once checked in",
	}

	opOrganize := Organize{
		Env: u,
		Metadata: organize_text.Metadata{
			RepoId: qg.RepoId,
			OptionCommentSet: organize_text.MakeOptionCommentSet(
				map[string]organize_text.OptionComment{
					"delete": flagDelete,
				},
				&organize_text.OptionCommentUnknown{
					Value: "instructions: to prevent an object from being checked in, delete it entirely",
				},
				organize_text.OptionCommentWithKey{
					Key:           "delete",
					OptionComment: flagDelete,
				},
			),
		},
		DontUseQueryGroupForOrganizeMetadata: true,
	}

	ui.Log().Print(qg)

	var organizeResults organize_text.OrganizeResults

	if organizeResults, err = opOrganize.RunWithQueryGroup(
		qg,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var changes organize_text.Changes

	if changes, err = organize_text.ChangesFromResults(
		u.GetConfig().PrintOptions,
		organizeResults,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = changes.After.Each(
		func(el sku.ExternalLike) (err error) {
			if err = u.GetStore().CreateOrUpdate(
				el,
				objekte_mode.ModeCreate,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			if !op.Delete {
				return
			}

			if err = u.GetStore().DeleteExternalLike(
				qg.RepoId,
				el,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Unlock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
