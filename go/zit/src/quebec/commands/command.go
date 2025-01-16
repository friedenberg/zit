package commands

import (
	"flag"
	"fmt"

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

func registerCommand(name string, cmd Command) {
	registerCommandOld(name, cmd)
}

func registerCommandOld(n string, commandOrCommandBuildFunc any) {
	f := flag.NewFlagSet(n, flag.ExitOnError)

	if _, ok := commands[n]; ok {
		panic("command added more than once: " + n)
	}

	switch cmd := commandOrCommandBuildFunc.(type) {
	case Command:
		wrapper := command.Wrapper{
			FlagSet: f,
			Command: cmd,
		}

		wrapper.SetFlagSet(f)

		commands[n] = wrapper

	case func(*flag.FlagSet) CommandOld:
		commands[n] = cmd(f)

	case func(*flag.FlagSet) WithLocalWorkingCopy:
		commands[n] = commandWithLocalWorkingCopy{
			Command: cmd(f),
			FlagSet: f,
		}

	default:
		panic(fmt.Sprintf("command or command build func not supported: %T", cmd))
	}
}

func registerCommandWithQuery(
	n string,
	makeFunc func(*flag.FlagSet) WithQuery,
) {
	registerCommandOld(
		n,
		func(f *flag.FlagSet) WithLocalWorkingCopy {
			cmd := &commandWithQuery{
				Command: makeFunc(f),
			}

			cmd.SetFlagSet(f)

			return cmd
		},
	)
}
