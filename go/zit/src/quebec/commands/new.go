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
	"code.linenisgreat.com/zit/go/zit/src/november/env"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

type New struct {
	ids.RepoId
	Delete bool
	Count  int
	// TODO combine organize and edit and refactor
	Organize  bool
	Edit      bool
	PrintOnly bool
	Filter    script_value.ScriptValue
	Shas      bool

	sku.Proto
}

func init() {
	registerCommand(
		"new",
		func(f *flag.FlagSet) Command {
			c := &New{}

			f.Var(&c.RepoId, "kasten", "none or Browser")

			f.BoolVar(
				&c.Shas,
				"shas",
				false,
				"treat arguments as blobs that are already checked in",
			)

			f.BoolVar(
				&c.Delete,
				"delete",
				false,
				"delete the zettel and blob after successful checkin",
			)

			f.BoolVar(
				&c.Organize,
				"organize",
				false,
				"open organize",
			)

			f.BoolVar(
				&c.Edit,
				"edit",
				true,
				"create a new empty zettel and open EDITOR or VISUAL for editing and then commit the resulting changes",
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

			return c
		},
	)
}

func (c New) ValidateFlagsAndArgs(
	u *env.Local,
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

func (c *New) Run(u *env.Local, args ...string) (err error) {
	if err = c.ValidateFlagsAndArgs(u, args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	cotfo := checkout_options.TextFormatterOptions{}

	f := object_metadata.MakeTextFormat(
		u.GetDirectoryLayout(),
		nil,
	)

	var zts sku.TransactedMutableSet

	if len(args) == 0 {
		emptyOp := user_ops.WriteNewZettels{
			Local: u,
		}

		if zts, err = emptyOp.RunMany(c.Proto, c.Count); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else if c.Shas {
		opCreateFromShas := user_ops.CreateFromShas{
			Local:   u,
			Proto: c.Proto,
		}

		if zts, err = opCreateFromShas.Run(args...); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		opCreateFromPath := user_ops.CreateFromPaths{
			Local:        u,
			TextParser: f,
			Filter:     c.Filter,
			Delete:     c.Delete,
			Proto:      c.Proto,
		}

		if zts, err = opCreateFromPath.Run(args...); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	// TODO make mutually exclusive with organize
	if c.Edit {
		opCheckout := user_ops.Checkout{
			Local: u,
			Options: checkout_options.Options{
				CheckoutMode: checkout_mode.MetadataAndBlob,
				OptionsWithoutMode: checkout_options.OptionsWithoutMode{
					TextFormatterOptions: cotfo,
				},
			},
			Edit: true,
		}

		if _, err = opCheckout.Run(zts); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if c.Organize {
		opOrganize := user_ops.Organize{
			Local: u,
		}

		if err = opOrganize.Metadata.SetFromObjectMetadata(
			&c.Metadata,
			ids.RepoId{},
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		var results organize_text.OrganizeResults

		if results, err = opOrganize.RunWithTransacted(nil, zts); err != nil {
			err = errors.Wrap(err)
			return
		}

		if _, err = u.LockAndCommitOrganizeResults(results); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
