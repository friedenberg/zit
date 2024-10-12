package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/script_value"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/organize_text"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

// fold into Checkin
type Add struct {
	AllowDupes         bool
	Delete             bool
	OpenBlob           bool
	CheckoutBlobAndRun string
	Organize           bool
	Filter             script_value.ScriptValue

	sku.Proto
}

func init() {
	registerCommandWithQuery(
		"add-old",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Add{}

			f.BoolVar(
				&c.AllowDupes,
				"allow-dupes",
				false,
				"permit added blobs to be duplicates (have the same exact content)",
			)

			f.BoolVar(
				&c.Delete,
				"delete",
				false,
				"delete the zettel and blob after successful checkin",
			)

			f.BoolVar(&c.OpenBlob, "open-blobs", false, "also open the blobs")

			f.StringVar(
				&c.CheckoutBlobAndRun,
				"each-blob",
				"",
				"checkout each Blob and run a utility",
			)

			f.BoolVar(&c.Organize, "organize", false, "")

			c.AddToFlagSet(f)

			ui.TodoP2(
				"add support for restricted query to specific gattung",
			)

			return c
		},
	)
}

func (c Add) ModifyBuilder(b *query.Builder) {
	b.WithDefaultGenres(ids.MakeGenre(genres.Zettel)).
		WithDoNotMatchEmpty()
}

func (c *Add) RunWithQuery(
	u *env.Env,
	qg *query.Group,
) (err error) {
	zettelsFromBlobOp := user_ops.ZettelFromExternalBlob{
		Env:        u,
		Proto:      c.Proto,
		Filter:     c.Filter,
		Delete:     c.Delete,
		AllowDupes: c.AllowDupes,
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var zettelsFromBlobResults sku.TransactedMutableSet

	if zettelsFromBlobResults, err = zettelsFromBlobOp.Run(qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Unlock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.openBlobIfNecessary(u, zettelsFromBlobResults); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !c.Organize {
		return
	}

	opOrganize := user_ops.Organize{
		Env: u,
	}

	if err = opOrganize.Metadata.SetFromObjectMetadata(
		&c.Metadata,
		qg.RepoId,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	opOrganize.Metadata.TagSet = u.GetConfig().DefaultTags.CloneSetPtrLike()

	var results organize_text.OrganizeResults

	if results, err = opOrganize.RunWithTransacted(
		nil,
		zettelsFromBlobResults,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = u.LockAndCommitOrganizeResults(results); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Add) openBlobIfNecessary(
	u *env.Env,
	zettels sku.TransactedSet,
) (err error) {
	if !c.OpenBlob && c.CheckoutBlobAndRun == "" {
		return
	}

	opCheckout := user_ops.Checkout{
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
