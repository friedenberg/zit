package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/script_value"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/organize_text"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

type New struct {
	ids.RepoId
	Count int
	// TODO combine organize and edit and refactor
	command_components.Checkout
	PrintOnly bool
	Filter    script_value.ScriptValue
	Shas      bool

	sku.Proto
}

func init() {
	registerCommand(
		"new",
		func(f *flag.FlagSet) CommandWithRepo {
			c := &New{}

			f.Var(&c.RepoId, "kasten", "none or Browser")

			f.BoolVar(
				&c.Shas,
				"shas",
				false,
				"treat arguments as blobs that are already checked in",
			)

			f.IntVar(
				&c.Count,
				"count",
				1,
				"when creating new empty zettels, how many to create. otherwise ignored",
			)

			f.Var(
				&c.Filter,
				"filter",
				"a script to run for each file to transform it the standard zettel format",
			)

			c.AddToFlagSet(f)
			c.Checkout.SetFlagSet(f)

			return c
		},
	)
}

func (c New) ValidateFlagsAndArgs(
	u *repo_local.Repo,
	args ...string,
) (err error) {
	if u.GetConfig().DryRun && len(args) == 0 {
		err = errors.Errorf(
			"when -dry-run is set, paths to existing zettels must be provided",
		)
		return
	}

	return
}

func (c *New) RunWithRepo(u *repo_local.Repo, args ...string) {
	if err := c.ValidateFlagsAndArgs(u, args...); err != nil {
		u.CancelWithError(err)
	}

	cotfo := checkout_options.TextFormatterOptions{}

	f := object_metadata.MakeTextFormat(
		object_metadata.Dependencies{
			DirLayout: u.GetRepoLayout().Layout,
			BlobStore: u.GetRepoLayout(),
		},
	)

	var zts sku.TransactedMutableSet

	if len(args) == 0 {
		emptyOp := user_ops.WriteNewZettels{
			Repo: u,
		}

		{
			var err error

			if zts, err = emptyOp.RunMany(c.Proto, c.Count); err != nil {
				u.CancelWithError(err)
			}
		}
	} else if c.Shas {
		opCreateFromShas := user_ops.CreateFromShas{
			Repo:  u,
			Proto: c.Proto,
		}

		{
			var err error

			if zts, err = opCreateFromShas.Run(args...); err != nil {
				u.CancelWithError(err)
			}
		}
	} else {
		opCreateFromPath := user_ops.CreateFromPaths{
			Repo:       u,
			TextParser: f,
			Filter:     c.Filter,
			Delete:     c.Delete,
			Proto:      c.Proto,
		}

		{
			var err error

			if zts, err = opCreateFromPath.Run(args...); err != nil {
				if errors.IsNotExist(err) {
					u.CancelWithBadRequestf("Expected a valid file path. Did you mean to add `-description`?")
				} else {
					u.CancelWithError(err)
				}
			}
		}
	}

	// TODO make mutually exclusive with organize
	if c.Edit {
		opCheckout := user_ops.Checkout{
			Repo: u,
			Options: checkout_options.Options{
				CheckoutMode: checkout_mode.MetadataAndBlob,
				OptionsWithoutMode: checkout_options.OptionsWithoutMode{
					Workspace: true,
					StoreSpecificOptions: store_fs.CheckoutOptions{
						TextFormatterOptions: cotfo,
					},
				},
			},
			Edit: true,
		}

		if _, err := opCheckout.Run(zts); err != nil {
			u.CancelWithError(err)
		}
	}

	if c.Organize {
		opOrganize := user_ops.Organize{
			Repo: u,
		}

		if err := opOrganize.Metadata.SetFromObjectMetadata(
			&c.Metadata,
			ids.RepoId{},
		); err != nil {
			u.CancelWithError(err)
		}

		var results organize_text.OrganizeResults

		{
			var err error

			if results, err = opOrganize.RunWithTransacted(nil, zts); err != nil {
				u.CancelWithError(err)
			}
		}

		if _, err := u.LockAndCommitOrganizeResults(results); err != nil {
			u.CancelWithError(err)
		}
	}
}
