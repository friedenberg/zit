package user_ops

import (
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/object_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
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
	u *env.Local,
	qg *query.Group,
) (err error) {
	var l sync.Mutex
	results := sku.MakeSkuTypeSetMutable()

	if op.Organize {
		if err = op.runOrganize(u, qg, results); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = u.GetStore().QuerySkuType(
			qg,
			func(co sku.SkuType) (err error) {
				l.Lock()
				defer l.Unlock()

				return results.Add(co.Clone())
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	processed := sku.MakeTransactedMutableSet()
	sortedResults := quiter.ElementsSorted(
		results,
		func(left, right sku.SkuType) bool {
			return left.String() < right.String()
		},
	)

	for _, co := range sortedResults {
		external := co.GetSkuExternal()

		if co.GetState() == checked_out_state.Untracked &&
			(co.GetSkuExternal().GetGenre() == genres.Zettel ||
				co.GetSkuExternal().GetGenre() == genres.Blob) {
			if external.Metadata.IsEmpty() {
				continue
			}

			if err = u.GetStore().UpdateTransactedFromBlobs(
				co,
			); err != nil {
				if errors.Is(err, external_store.ErrUnsupportedOperation{}) {
					err = nil
				} else {
					err = errors.Wrap(err)
					return
				}
			}

			external.ObjectId.Reset()

			if err = u.GetStore().CreateOrUpdate(
				external,
				object_mode.ModeApplyProto,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			if op.Proto.Apply(external, genres.Zettel) {
				if err = u.GetStore().CreateOrUpdate(
					external.GetSku(),
					object_mode.ModeEmpty,
				); err != nil {
					err = errors.Wrap(err)
					return
				}
			}
		} else {
			if err = u.GetStore().CreateOrUpdateCheckedOut(
				co,
				!op.Delete,
			); err != nil {
				err = errors.Wrapf(err, "CheckedOut: %s", co)
				return
			}
		}

		if !op.Delete {
			continue
		}

		if err = u.GetStore().DeleteCheckedOut(co); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = processed.Add(co.GetSkuExternal().CloneTransacted()); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = u.Unlock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = op.openBlobIfNecessary(u, processed); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (op Checkin) runOrganize(
	u *env.Local,
	qg *query.Group,
	results sku.SkuTypeSetMutable,
) (err error) {
	flagDelete := organize_text.OptionCommentBooleanFlag{
		Value:   &op.Delete,
		Comment: "delete once checked in",
	}

	opOrganize := Organize{
		Local: u,
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

	// TODO switch to using SkuType?
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

	if err = changes.After.Each(
		func(co sku.SkuType) (err error) {
			if err = results.Add(co.Clone()); err != nil {
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

func (c Checkin) openBlobIfNecessary(
	u *env.Local,
	objects sku.TransactedSet,
) (err error) {
	if !c.OpenBlob && c.CheckoutBlobAndRun == "" {
		return
	}

	opCheckout := Checkout{
		Local: u,
		Options: checkout_options.Options{
			CheckoutMode: checkout_mode.BlobOnly,
		},
		Utility: c.CheckoutBlobAndRun,
	}

	if _, err = opCheckout.Run(objects); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
