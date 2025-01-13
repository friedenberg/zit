package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/hotel/repo_layout"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
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
			c.Config.RepoType = repo_type.TypeReadWrite

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
	repo := cmd.OnTheFirstDay(
		dependencies.Context,
		dependencies.Config,
		env.Options{},
	)

	remote := cmd.MakeRemote(repo.GetEnv(), cmd.GetFlagSet().Args()[0])

	qg := cmd.MakeQueryGroup(
		query.MakeBuilderOptions(cmd),
		repo,
		cmd.Args()[1:]...,
	)

	if err := repo.PullQueryGroupFromRemote(
		remote,
		qg,
		cmd.RemoteTransferOptions.WithPrintCopies(true),
	); err != nil {
		dependencies.CancelWithError(err)
	}
}
