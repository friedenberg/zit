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

func (op Checkin) Run(
	u *env.Env,
	qg *query.Group,
) (err error) {
	if op.Organize {
		if qg, err = op.runOrganize(u, qg); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = op.runNoOrganize(u, qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (op Checkin) runOrganize(
	u *env.Env,
	qgOriginal *query.Group,
) (qgModified *query.Group, err error) {
	flagDelete := organize_text.OptionCommentBooleanFlag{
		Value:   &op.Delete,
		Comment: "delete once checked in",
	}

	opOrganize := Organize{
		Env: u,
		Metadata: organize_text.Metadata{
			OptionCommentSet: organize_text.MakeOptionCommentSet(
				map[string]organize_text.OptionComment{
					"delete": flagDelete,
				},
				organize_text.OptionCommentUnknown(
					"instructions: to prevent an object from being checked in, delete it entirely",
				),
				organize_text.OptionCommentWithKey{
					Key:           "delete",
					OptionComment: flagDelete,
				},
			),
		},
		DontUseQueryGroupForOrganizeMetadata: true,
	}

	ui.Log().Print(qgOriginal)

	var organizeResults organize_text.OrganizeResults

	if organizeResults, err = opOrganize.RunWithQueryGroup(
		qgOriginal,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	u.Lock()
	defer errors.Deferred(&err, u.Unlock)

	if qgModified, _, err = u.QueryGroupFromRemainingOrganizeResults(
		organizeResults,
		qgOriginal.RepoId,
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

	return
}
