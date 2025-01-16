package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

type Init struct {
	*flag.FlagSet
	command_components.Genesis
}

func init() {
	registerCommand(
		"init",
		func(f *flag.FlagSet) Command {
			c := &Init{}

			c.SetFlagSet(f)

			return c
		},
	)
}

func (cmd *Init) SetFlagSet(f *flag.FlagSet) {
	cmd.FlagSet = f
	cmd.Genesis.SetFlagSet(f)
}

func (c *Init) GetFlagSet() *flag.FlagSet {
	return c.FlagSet
}

func (cmd *Init) Run(
	dependencies command.Dep,
) {
	cmd.OnTheFirstDay(
		dependencies.Context,
		dependencies.Config,
		env.Options{},
	)
}
