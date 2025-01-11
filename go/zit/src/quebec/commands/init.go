package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type Init struct {
	*flag.FlagSet
	repo_local.BigBang
}

func init() {
	registerCommand(
		"init",
		func(f *flag.FlagSet) Command {
			c := &Init{
				FlagSet: f,
				BigBang: repo_local.BigBang{
					Config: immutable_config.Default(),
				},
			}

			c.BigBang.SetFlagSet(f)

			return c
		},
	)
}

func (c Init) GetFlagSet() *flag.FlagSet {
	return c.FlagSet
}

func (c Init) Run(
	dependencies Dependencies,
) {
	var repo *repo_local.Repo

	{
		var err error

		if repo, err = c.BigBang.Start(
			dependencies.Context,
			dependencies.Config,
			env.Options{},
		); err != nil {
			repo.CancelWithError(err)
		}
	}
}
