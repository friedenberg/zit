package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/script_value"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

type Add struct {
	Dedupe             bool
	Delete             bool
	OpenBlob           bool
	CheckoutBlobAndRun string
	Organize           bool
	Filter             script_value.ScriptValue

	sku.Proto
}

func init() {
	registerCommand(
		"add",
		func(f *flag.FlagSet) Command {
			c := &Add{}

			f.BoolVar(
				&c.Dedupe,
				"dedupe",
				false,
				"deduplicate added Zettels based on blob sha",
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

func (c Add) Run(
	u *env.Env,
	args ...string,
) (err error) {
	zettelsFromBlobOp := user_ops.ZettelFromExternalBlob{
		Env:    u,
		Proto:  c.Proto,
		Filter: c.Filter,
		Delete: c.Delete,
		Dedupe: c.Dedupe,
	}

	var zettelsFromBlobResults sku.TransactedMutableSet

	fds := fd.MakeMutableSet()

	for _, v := range args {
		if v == "." {
			if err = u.GetStore().GetCwdFiles().GetBlobFDs().Each(fds.Add); err != nil {
				err = errors.Wrap(err)
				return
			}

			break
		} else if v == "" {
			continue
		}

		var f fd.FD

		if err = f.Set(v); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = fds.Add(&f); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if zettelsFromBlobResults, err = zettelsFromBlobOp.Run(fds); err != nil {
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
		Env:      u,
		Metadata: c.Metadata,
	}

	if err = u.GetConfig().DefaultTags.EachPtr(
		opOrganize.Metadata.AddTagPtr,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = opOrganize.Run(
		nil,
		zettelsFromBlobResults,
	); err != nil {
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
