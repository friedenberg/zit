package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type (
	CommandOld interface {
		GetFlagSet() *flag.FlagSet
		Run(command.Dep)
	}

	Command = command.Command
)

type CompleteWithRepo interface {
	Complete(u *local_working_copy.Repo, args ...string)
}

var commands = map[string]CommandOld{}

func Commands() map[string]CommandOld {
	return commands
}

func registerCommand(n string, cmd Command) {
	f := flag.NewFlagSet(n, flag.ExitOnError)

	if _, ok := commands[n]; ok {
		panic("command added more than once: " + n)
	}

	wrapper := command.Wrapper{
		FlagSet: f,
		Command: cmd,
	}

	wrapper.SetFlagSet(f)

	commands[n] = wrapper
}
