package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type Clone struct {
	*flag.FlagSet
	repo_local.BigBang
	ComponentRemote
	ComponentQuery
}

// TODO transition this to CommandWithDependencies instead of
// CommandWithRemoteAndQuery
func init() {
	registerCommand(
		"clone",
		func(f *flag.FlagSet) CommandWithDependencies {
			c := &Clone{
				BigBang: repo_local.BigBang{
					Config:             immutable_config.Default(),
					ExcludeDefaultType: true,
				},
			}

			c.SetFlagSet(f)

			return c
		},
	)
}

func (cmd *Clone) GetCommandWithDependencies() CommandWithDependencies {
	return cmd
}

func (cmd *Clone) GetFlagSet() *flag.FlagSet {
	return cmd.FlagSet
}

func (cmd *Clone) SetFlagSet(f *flag.FlagSet) {
	cmd.BigBang.AddToFlagSet(f)
	cmd.ComponentRemote.SetFlagSet(f)
	cmd.ComponentQuery.SetFlagSet(f)
}

func (c Clone) DefaultSigil() ids.Sigil {
	return ids.MakeSigil(ids.SigilHistory, ids.SigilHidden)
}

func (c Clone) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.InventoryList)
	// return ids.MakeGenre(genres.TrueGenre()...)
}

func (c Clone) RunWithDependencies(
	dependencies Dependencies,
) {
	var local *repo_local.Repo

	{
		var err error

		if local, err = c.BigBang.Start(
			dependencies.Context,
			dependencies.Config,
		); err != nil {
			local.CancelWithError(err)
		}

		defer dependencies.MustWithContext(local.GetDirLayout().ResetTempOnExit)
		defer local.MustFlush(local)
	}

	remote := c.MakeRemote(local.Env, c.GetFlagSet().Args()[0])

	var qg *query.Group

	if err := local.PullQueryGroupFromRemote(
		remote,
		qg,
		true,
	); err != nil {
		local.CancelWithError(err)
	}
}
