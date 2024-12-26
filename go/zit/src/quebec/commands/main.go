package commands

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
)

// TODO switch to returning result
func Run(
	ctx errors.Context,
	args ...string,
) (exitStatus int) {
	var cmd CommandWithDependencies

	if len(os.Args) < 1 {
		ui.Log().Print("printing usage")
		return PrintUsage(nil)
	}

	if len(os.Args) == 1 {
		return PrintUsage(errors.Errorf("No subcommand provided."))
	}

	cmds := Commands()
	specifiedSubcommand := os.Args[1]

	ok := false

	if cmd, ok = cmds[specifiedSubcommand]; !ok {
		return PrintUsage(
			errors.BadRequestf("No subcommand '%s'", specifiedSubcommand),
		)
	}

	args = os.Args[2:]

	// TODO customize command flag parsing and env / dir layout creation based on
	// type of command
	configCli := config_mutable_cli.Default()
	configCli.AddToFlags(cmd.GetFlagSet())

	if err := cmd.GetFlagSet().Parse(args); err != nil {
		ctx.CancelWithError(err)
		return
	}

	exitStatus = cmd.RunWithDependencies(
		Dependencies{
			Context: ctx,
			Config:  configCli,
		},
	)

	return
}
