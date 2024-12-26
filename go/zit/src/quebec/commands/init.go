package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/delta/debug"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
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
		func(f *flag.FlagSet) CommandWithDependencies {
			c := &Init{
				FlagSet: f,
				BigBang: repo_local.BigBang{
					Config: immutable_config.Default(),
				},
			}

			c.BigBang.AddToFlagSet(f)

			return c
		},
	)
}

func (c Init) GetFlagSet() *flag.FlagSet {
	return c.FlagSet
}

func (c Init) RunWithDependencies(
	dependencies Dependencies,
) {
	if c.OverrideXDGWithCwd {
		// TODO init dirLayout with XDG override
	}

	// TODO use options when making dirLayout
	var dirLayout dir_layout.Layout

	{
		var err error

		if dirLayout, err = dir_layout.MakeDefault(
			dependencies.Debug,
		); err != nil {
			dependencies.CancelWithError(err)
			return
		}
	}

	// TODO move to env
	if _, err := debug.MakeContext(
		dependencies.Context,
		dependencies.Debug,
	); err != nil {
		dependencies.CancelWithError(err)
		return
	}

	env := env.Make(
		dependencies.Context,
		c.GetFlagSet(),
		dependencies.Config,
		dirLayout,
	)

	defer func() {
		if err := env.GetDirLayout().ResetTempOnExit(
			dependencies.Context,
		); err != nil {
			dependencies.CancelWithError(err)
			return
		}
	}()

	var repo *repo_local.Repo

	{
		var err error

		if repo, err = repo_local.Make(
			env,
			repo_local.OptionsEmpty,
		); err != nil {
			env.CancelWithError(err)
			return
		}

		defer env.Flusher(repo)
	}

	if err := repo.Start(c.BigBang); err != nil {
		env.CancelWithError(err)
		return
	}
}
