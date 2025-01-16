package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
)

func Run(ctx *errors.Context, args ...string) {
	if len(args) <= 1 {
		PrintUsage(
			ctx,
			errors.BadRequestf("No subcommand provided."),
		)
	}

	cmds := command.Commands()

	var cmd command.Command
	var ok bool

	specifiedSubcommand := args[1]

	if cmd, ok = cmds[specifiedSubcommand]; !ok {
		PrintUsage(
			ctx,
			errors.BadRequestf("No subcommand '%s'", specifiedSubcommand),
		)
	}

	f := flag.NewFlagSet(specifiedSubcommand, flag.ContinueOnError)
	cmd.SetFlagSet(f)

	args = args[2:]

	configCli := config_mutable_cli.Default()
	configCli.SetFlagSet(f)

	if err := f.Parse(args); err != nil {
		ctx.CancelWithError(err)
	}

	dep := command.Dep{
		Context: ctx,
		Config:  configCli,
		FlagSet: f,
	}

	cmd.Run(dep)
}
