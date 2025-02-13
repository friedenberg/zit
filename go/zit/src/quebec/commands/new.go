package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/script_value"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/organize_text"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

func init() {
	command.Register("new", &New{})
}

type New struct {
	command_components.LocalWorkingCopy

	complete command_components.Complete

	ids.RepoId
	Count int
	// TODO combine organize and edit and refactor
	command_components.Checkout
	PrintOnly bool
	Filter    script_value.ScriptValue
	Shas      bool

	sku.Proto
}

func (cmd *New) SetFlagSet(f *flag.FlagSet) {
	f.Var(&cmd.RepoId, "kasten", "none or Browser")

	f.BoolVar(
		&cmd.Shas,
		"shas",
		false,
		"treat arguments as blobs that are already checked in",
	)

	f.IntVar(
		&cmd.Count,
		"count",
		1,
		"when creating new empty zettels, how many to create. otherwise ignored",
	)

	f.Var(
		&cmd.Filter,
		"filter",
		"a script to run for each file to transform it the standard zettel format",
	)

	cmd.complete.SetFlagsProto(
		&cmd.Proto,
		f,
		"description to use for new zettels",
		"tags added for new zettels",
		"type used for new zettels",
	)

	cmd.Checkout.SetFlagSet(f)
}

func (c New) ValidateFlagsAndArgs(
	u *local_working_copy.Repo,
	args ...string,
) (err error) {
	if u.GetConfig().GetCLIConfig().IsDryRun() && len(args) == 0 {
		err = errors.Errorf(
			"when -dry-run is set, paths to existing zettels must be provided",
		)
		return
	}

	return
}

func (cmd *New) Run(req command.Request) {
	args := req.Args()
	repo := cmd.MakeLocalWorkingCopy(req)

	if err := cmd.ValidateFlagsAndArgs(repo, args...); err != nil {
		repo.CancelWithError(err)
	}

	cotfo := checkout_options.TextFormatterOptions{}

	f := object_metadata.MakeTextFormat(
		object_metadata.Dependencies{
			EnvDir:    repo.GetEnvRepo(),
			BlobStore: repo.GetEnvRepo(),
		},
	)

	var zts sku.TransactedMutableSet

	if len(args) == 0 {
		emptyOp := user_ops.WriteNewZettels{
			Repo: repo,
		}

		{
			var err error

			if zts, err = emptyOp.RunMany(cmd.Proto, cmd.Count); err != nil {
				repo.CancelWithError(err)
			}
		}
	} else if cmd.Shas {
		opCreateFromShas := user_ops.CreateFromShas{
			Repo:  repo,
			Proto: cmd.Proto,
		}

		{
			var err error

			if zts, err = opCreateFromShas.Run(args...); err != nil {
				repo.CancelWithError(err)
			}
		}
	} else {
		opCreateFromPath := user_ops.CreateFromPaths{
			Repo:       repo,
			TextParser: f,
			Filter:     cmd.Filter,
			Delete:     cmd.Delete,
			Proto:      cmd.Proto,
		}

		{
			var err error

			if zts, err = opCreateFromPath.Run(args...); err != nil {
				if errors.IsNotExist(err) {
					repo.CancelWithBadRequestf("Expected a valid file path. Did you mean to add `-description`?")
				} else {
					repo.CancelWithError(err)
				}
			}
		}
	}

	envWorkspace := repo.GetEnvWorkspace()

	// TODO make mutually exclusive with organize
	if cmd.Edit {
		opCheckout := user_ops.Checkout{
			Repo: repo,
			Options: checkout_options.Options{
				CheckoutMode: checkout_mode.MetadataAndBlob,
				OptionsWithoutMode: checkout_options.OptionsWithoutMode{
					IgnoreWorkspace: !envWorkspace.InWorkspace(),
					StoreSpecificOptions: store_fs.CheckoutOptions{
						TextFormatterOptions: cotfo,
					},
				},
			},
			Edit: true,
		}

		if _, err := opCheckout.Run(zts); err != nil {
			repo.CancelWithError(err)
		}
	}

	if cmd.Organize {
		opOrganize := user_ops.Organize{
			Repo: repo,
		}

		if err := opOrganize.Metadata.SetFromObjectMetadata(
			&cmd.Metadata,
			ids.RepoId{},
		); err != nil {
			repo.CancelWithError(err)
		}

		var results organize_text.OrganizeResults

		{
			var err error

			if results, err = opOrganize.RunWithTransacted(nil, zts); err != nil {
				repo.CancelWithError(err)
			}
		}

		if _, err := repo.LockAndCommitOrganizeResults(results); err != nil {
			repo.CancelWithError(err)
		}
	}
}
