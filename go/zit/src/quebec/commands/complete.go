package commands

import (
	"flag"
	"io"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_local"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	command.Register(
		"complete",
		&Complete{},
	)
}

type Complete struct {
	command_components.Env
	bashStyle  bool
	inProgress string
}

func (cmd Complete) GetDescription() command.Description {
	return command.Description{
		Short: "complete a command-line",
	}
}

func (cmd *Complete) SetFlagSet(f *flag.FlagSet) {
	f.BoolVar(&cmd.bashStyle, "bash-style", false, "")
	f.StringVar(&cmd.inProgress, "in-progress", "", "")
}

func (cmd Complete) Run(req command.Request) {
	cmds := command.Commands()
	envLocal := cmd.MakeEnv(req)

	commandLine := command.CommandLine{
		Args:       req.Args(),
		InProgress: cmd.inProgress,
	}

	// TODO determine state:
	// bare: `zit`
	// subcommand or arg or flag:
	//  - `zit subcommand`
	//  - `zit subcommand -flag=true`
	//  - `zit subcommand -flag value`
	// flag: `zit subcommand -flag`
	lastArg, hasLastArg := commandLine.LastArg()

	if !hasLastArg {
		cmd.completeSubcommands(envLocal, commandLine, cmds)
		return
	}

	name := req.Args()[0]
	subcmd, foundSubcmd := cmds[name]

	if !foundSubcmd {
		cmd.completeSubcommands(envLocal, commandLine, cmds)
		return
	}

	flagSet := flag.NewFlagSet(name, flag.ContinueOnError)
	flagSet.SetOutput(io.Discard)
	subcmd.SetFlagSet(flagSet)

	if cmd.completeSubcommandFlags(
		req,
		envLocal,
		subcmd,
		flagSet,
		commandLine,
		lastArg,
	) {
		return
	}

	cmd.completeSubcommandArgs(envLocal, subcmd, commandLine)
}

func (cmd Complete) completeSubcommands(
	envLocal env_local.Env,
	commandLine command.CommandLine,
	cmds map[string]command.Command,
) {
	for name, subcmd := range cmds {
		cmd.completeSubcommand(envLocal, name, subcmd)
	}
}

func (cmd Complete) completeSubcommand(
	envLocal env_local.Env,
	name string,
	subcmd command.Command,
) {
	var shortDescription string

	if hasDescription, ok := subcmd.(command.HasDescription); ok {
		description := hasDescription.GetDescription()
		shortDescription = description.Short
	}

	if shortDescription != "" {
		envLocal.GetUI().Printf("%s\t%s", name, shortDescription)
	} else {
		envLocal.GetUI().Printf("%s", name)
	}
}

func (cmd Complete) completeSubcommandArgs(
	envLocal env_local.Env,
	subcmd command.Command,
	commandLine command.CommandLine,
) {
	if subcmd == nil {
		return
	}

	completer, isCompleter := subcmd.(command.Completer)

	if !isCompleter {
		return
	}

	completer.Complete(envLocal, commandLine)
}

func (cmd Complete) completeSubcommandFlags(
	req command.Request,
	envLocal env_local.Env,
	subcmd command.Command,
	flagSet *flag.FlagSet,
	commandLine command.CommandLine,
	lastArg string,
) (flagComplete bool) {
	if subcmd == nil {
		return
	}

	if strings.HasPrefix(lastArg, "-") {
		flagComplete = true
	} else if commandLine.InProgress != "" && len(commandLine.Args) > 1 {
		lastArg = commandLine.Args[len(commandLine.Args)-2]
		commandLine.InProgress = ""
		flagComplete = strings.HasPrefix(lastArg, "-")
	}

	if commandLine.InProgress != "" {
		flagSet.VisitAll(func(flag *flag.Flag) {
			envLocal.GetUI().Printf("-%s\t%s", flag.Name, flag.Usage)
		})
	} else if err := flagSet.Parse([]string{lastArg}); err != nil {
		cmd.completeSubcommandFlagOnParseError(
			req,
			envLocal,
			subcmd,
			flagSet,
			err,
		)
	} else {
		flagSet.VisitAll(func(flag *flag.Flag) {
			envLocal.GetUI().Printf("-%s\t%s", flag.Name, flag.Usage)
		})
	}

	return
}

func (cmd Complete) completeSubcommandFlagOnParseError(
	req command.Request,
	envLocal env_local.Env,
	subcmd command.Command,
	flagSet *flag.FlagSet,
	err error,
) {
	if subcmd == nil {
		return
	}

	after, found := strings.CutPrefix(
		err.Error(),
		"flag needs an argument: -",
	)

	if !found {
		envLocal.CancelWithError(err)
		return
	}

	switch after {
	case "tags":
		local := command_components.LocalWorkingCopy{}.MakeLocalWorkingCopy(req)
		command_components.QueryGroup{}.CompleteTagsWithRepo(
			req,
			local,
		)

	default:
		ui.Debug().Print("to complete:", after)
	}
}
