package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/hotel/repo_layout"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

type Clone struct {
	*flag.FlagSet
	command_components.Genesis
	command_components.RemoteTransfer
	command_components.QueryGroup
}

func init() {
	registerCommand(
		"clone",
		func(f *flag.FlagSet) Command {
			c := &Clone{
				FlagSet: f,
				Genesis: command_components.Genesis{
					BigBang: repo_layout.BigBang{
						ExcludeDefaultType: true,
					},
				},
			}

			c.SetFlagSet(f)
			c.Config.RepoType = repo_type.TypeWorkingCopy

			return c
		},
	)
}

func (cmd *Clone) GetCommandWithDependencies() Command {
	return cmd
}

func (cmd *Clone) GetFlagSet() *flag.FlagSet {
	return cmd.FlagSet
}

func (cmd *Clone) SetFlagSet(f *flag.FlagSet) {
	cmd.Genesis.SetFlagSet(f)
	cmd.RemoteTransfer.SetFlagSet(f)
	cmd.QueryGroup.SetFlagSet(f)
}

func (c Clone) DefaultSigil() ids.Sigil {
	return ids.MakeSigil(ids.SigilHistory, ids.SigilHidden)
}

func (c Clone) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.InventoryList)
	// return ids.MakeGenre(genres.TrueGenre()...)
}

func (cmd Clone) Run(
	dependencies Dependencies,
) {
	repoGeneric := cmd.OnTheFirstDay(
		dependencies.Context,
		dependencies.Config,
		env.Options{},
	)

	remote := cmd.MakeRemote(
		repoGeneric.GetRepoLayout().GetEnv(),
		cmd.GetFlagSet().Args()[0],
	)

	switch local := repoGeneric.(type) {
	default:
		dependencies.CancelWithBadRequestf(
			"unsupported repo type: %q (%T)",
			local.GetRepoLayout().GetConfig().GetRepoType(),
			local,
		)

	case repo.WorkingCopy:
		qg := cmd.MakeQueryGroup(
			query.MakeBuilderOptions(cmd),
			local,
			cmd.Args()[1:]...,
		)

		if err := local.PullQueryGroupFromRemote(
			remote,
			qg,
			cmd.RemoteTransferOptions.WithPrintCopies(true),
		); err != nil {
			dependencies.CancelWithError(err)
		}

	case repo.Archive:
		remoteInventoryListStore := remote.GetInventoryListStore()
		localInventoryListStore := local.GetInventoryListStore()

		if err := remoteInventoryListStore.ReadAllInventoryLists(
			func(sk *sku.Transacted) (err error) {
				if err = localInventoryListStore.ImportInventoryList(
					remote.GetBlobStore(),
					sk,
				); err != nil {
					err = errors.Wrap(err)
					return
				}

				return
			},
		); err != nil {
			dependencies.CancelWithError(err)
		}
	}
}
