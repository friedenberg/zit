package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/echo/repo_layout"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type InitRelay struct {
	*flag.FlagSet
	repo_layout.BigBang
}

func init() {
	registerCommand(
		"init-relay",
		func(f *flag.FlagSet) Command {
			c := &InitRelay{
				FlagSet: f,
				BigBang: repo_layout.BigBang{},
			}

			c.SetFlagSet(f)
			c.Config.RepoType = repo_type.TypeRelay

			return c
		},
	)
}

func (c InitRelay) GetFlagSet() *flag.FlagSet {
	return c.FlagSet
}

func (c InitRelay) Run(
	dependencies Dependencies,
) {
	repo_local.Genesis(
		c.BigBang,
		dependencies.Context,
		dependencies.Config,
		env.Options{},
	)
}
