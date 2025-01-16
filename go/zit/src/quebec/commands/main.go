package commands

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
)

func Run(
	ctx *errors.Context,
	args ...string,
) {
	if len(args) <= 1 {
		PrintUsage(
			ctx,
			errors.BadRequestf("No subcommand provided."),
		)
	}

	cmds := Commands()
	var cmd Command
	var ok bool

	specifiedSubcommand := args[1]

	if cmd, ok = cmds[specifiedSubcommand]; !ok {
		PrintUsage(
			ctx,
			errors.BadRequestf("No subcommand '%s'", specifiedSubcommand),
		)
	}

	args = args[2:]

	configCli := config_mutable_cli.Default()
	configCli.AddToFlags(cmd.GetFlagSet())

	if err := cmd.GetFlagSet().Parse(args); err != nil {
		ctx.CancelWithError(err)
	}

	cmd.Run(
		command.Dep{
			Context: ctx,
			Config:  configCli,
			FlagSet: cmd.GetFlagSet(),
		},
	)
}
