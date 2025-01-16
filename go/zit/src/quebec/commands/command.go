package commands

import (
	"flag"
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type Command interface {
	GetFlagSet() *flag.FlagSet
	Run(command.Dep)
}

type Command2 interface {
	Run(command.Dep)
	interfaces.CommandComponent
}

type commandWrapper struct {
	*flag.FlagSet
	Command2
}

func (wrapper commandWrapper) GetFlagSet() *flag.FlagSet {
	return wrapper.FlagSet
}

func (wrapper commandWrapper) SetFlagSet(f *flag.FlagSet) {
	wrapper.Command2.SetFlagSet(f)
}

type CompleteWithRepo interface {
	Complete(u *local_working_copy.Repo, args ...string)
}

var commands = map[string]Command{}

func Commands() map[string]Command {
	return commands
}

func registerCommand(
	n string,
	commandOrCommandBuildFunc any,
) {
	f := flag.NewFlagSet(n, flag.ExitOnError)

	if _, ok := commands[n]; ok {
		panic("command added more than once: " + n)
	}

	switch cmd := commandOrCommandBuildFunc.(type) {
	case Command2:
		wrapper := commandWrapper{
			FlagSet:  f,
			Command2: cmd,
		}

		wrapper.SetFlagSet(f)

		commands[n] = wrapper

	case func(*flag.FlagSet) Command:
		commands[n] = cmd(f)

	case func(*flag.FlagSet) WithLocalWorkingCopy:
		commands[n] = commandWithLocalWorkingCopy{
			Command: cmd(f),
			FlagSet: f,
		}

	case func(*flag.FlagSet) WithBlobStore:
		commands[n] = commandWithBlobStore{
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
	registerCommand(
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
