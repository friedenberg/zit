package user_ops

import (
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/organize_text"
	"code.linenisgreat.com/zit/go/zit/src/november/read_write_repo_local"
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
	u *read_write_repo_local.Repo,
	qg *query.Group,
) (err error) {
	var l sync.Mutex

	results := sku.MakeSkuTypeSetMutable()

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

	if op.Organize {
		if err = op.runOrganize(u, qg, results); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	var processed sku.TransactedMutableSet

	if processed, err = u.Checkin(
		results,
		op.Proto,
		op.Delete,
	); err != nil {
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
	u *read_write_repo_local.Repo,
	qg *query.Group,
	results sku.SkuTypeSetMutable,
) (err error) {
	flagDelete := organize_text.OptionCommentBooleanFlag{
		Value:   &op.Delete,
		Comment: "delete once checked in",
	}

	opOrganize := Organize{
		Repo: u,
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

	for _, co := range changes.After.AllSkuAndIndex() {
		if err = results.Add(co.Clone()); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (c Checkin) openBlobIfNecessary(
	u *read_write_repo_local.Repo,
	objects sku.TransactedSet,
) (err error) {
	if !c.OpenBlob && c.CheckoutBlobAndRun == "" {
		return
	}

	opCheckout := Checkout{
		Repo: u,
		Options: checkout_options.Options{
			CheckoutMode: checkout_mode.BlobOnly,
		},
		Utility: c.CheckoutBlobAndRun,
	}

	opCheckout.Workspace = true

	if _, err = opCheckout.Run(objects); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
