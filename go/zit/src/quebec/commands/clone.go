package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/echo/repo_layout"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

type Clone struct {
	*flag.FlagSet
	repo_layout.BigBang
	command_components.RemoteTransfer
	command_components.QueryGroup
}

func init() {
	registerCommand(
		"clone",
		func(f *flag.FlagSet) Command {
			c := &Clone{
				BigBang: repo_layout.BigBang{
					Config:             immutable_config.Default(),
					ExcludeDefaultType: true,
				},
			}

			c.SetFlagSet(f)

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
	cmd.FlagSet = f
	cmd.BigBang.SetFlagSet(f)
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

func (c Clone) Run(
	dependencies Dependencies,
) {
	var local *repo_local.Repo

	{
		var err error

		if local, err = repo_local.Genesis(
			c.BigBang,
			dependencies.Context,
			dependencies.Config,
			env.Options{},
		); err != nil {
			local.CancelWithError(err)
		}
	}

	remote := c.MakeRemote(local.Env, c.GetFlagSet().Args()[0])

	qg := c.MakeQueryGroup(c, local, c.Args()[1:]...)

	if err := local.PullQueryGroupFromRemote(
		remote,
		qg,
		c.RemoteTransferOptions.WithPrintCopies(true),
	); err != nil {
		local.CancelWithError(err)
	}
}
