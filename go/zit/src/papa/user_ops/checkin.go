package user_ops

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/object_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/organize_text"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type Checkin struct {
	Proto sku.Proto

	// TODO make flag family disambiguate these options
	// and use with other commands too
	Delete             bool
	Organize           bool
	CheckoutBlobAndRun string
	OpenBlob           bool
	Edit               bool // TODO add support back for this
}

func (op Checkin) Run(
	u *env.Env,
	qg *query.Group,
) (err error) {
	// TODO make organize use results in order to support open blob via organize
	// path
	results := sku.MakeTransactedMutableSet()

	if op.Organize {
		if err = op.runOrganize(u, qg, results); err != nil {
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
				cofs := col.(*sku.CheckedOut)
				z := col.GetSkuExternalLike().GetSku()

				if cofs.State == checked_out_state.Untracked &&
					(cofs.External.GetGenre() == genres.Zettel ||
						cofs.External.GetGenre() == genres.Blob) {
					if z.Metadata.IsEmpty() {
						return
					}

					if err = u.GetStore().GetStoreFS().UpdateTransactedFromBlobs(
						&cofs.External,
					); err != nil {
						err = errors.Wrap(err)
						return
					}

					z.ObjectId.Reset()

					if err = u.GetStore().CreateOrUpdate(
						z,
						object_mode.ModeApplyProto,
					); err != nil {
						err = errors.Wrap(err)
						return
					}

					if op.Proto.Apply(z, genres.Zettel) {
						if err = u.GetStore().CreateOrUpdate(
							z.GetSku(),
							object_mode.ModeEmpty,
						); err != nil {
							err = errors.Wrap(err)
							return
						}
					}
				} else {
					if err = u.GetStore().CreateOrUpdateCheckedOut(
						col,
						!op.Delete,
					); err != nil {
						err = errors.Wrap(err)
						return
					}
				}

				if !op.Delete {
					return
				}

				if err = u.GetStore().DeleteCheckedOutLike(col); err != nil {
					err = errors.Wrap(err)
					return
				}

				if err = results.Add(&cofs.External); err != nil {
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

	if err = op.openBlobIfNecessary(u, results); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (op Checkin) runOrganize(
	u *env.Env,
	qg *query.Group,
	results sku.TransactedMutableSet,
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
				object_mode.ModeCreate,
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

func (c Checkin) openBlobIfNecessary(
	u *env.Env,
	zettels sku.TransactedSet,
) (err error) {
	if !c.OpenBlob && c.CheckoutBlobAndRun == "" {
		return
	}

	opCheckout := Checkout{
		Env: u,
		Options: checkout_options.Options{
			CheckoutMode: checkout_mode.BlobOnly,
		},
		Utility: c.CheckoutBlobAndRun,
	}

	if _, err = opCheckout.Run(
		zettels,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
