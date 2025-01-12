package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/echo/repo_layout"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type Init struct {
	*flag.FlagSet
	repo_layout.BigBang
}

func init() {
	registerCommand(
		"init",
		func(f *flag.FlagSet) Command {
			c := &Init{
				FlagSet: f,
				BigBang: repo_layout.BigBang{},
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
	repo_local.Genesis(
		c.BigBang,
		dependencies.Context,
		dependencies.Config,
		env.Options{},
	)
}
