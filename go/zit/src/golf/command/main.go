package command

import (
	"flag"
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
)

type Command interface {
	GetFlagSet() *flag.FlagSet
	Run(Dep)
}

type Command2 interface {
	Run(Dep)
	interfaces.CommandComponent
}

type Dep struct {
	*errors.Context
	config_mutable_cli.Config
	*flag.FlagSet
}

var commands = map[string]Command{}

func Commands() map[string]Command {
	return commands
}

func Register(
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

	// case func(*flag.FlagSet) WithLocalWorkingCopy:
	// 	commands[n] = commandWithLocalWorkingCopy{
	// 		Command: cmd(f),
	// 		FlagSet: f,
	// 	}

	// case func(*flag.FlagSet) WithBlobStore:
	// 	commands[n] = commandWithBlobStore{
	// 		Command: cmd(f),
	// 		FlagSet: f,
	// 	}

	default:
		panic(fmt.Sprintf("command or command build func not supported: %T", cmd))
	}
}
