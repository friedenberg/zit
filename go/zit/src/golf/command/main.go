package command

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
)

type Dep struct {
	*errors.Context
	config_mutable_cli.Config
	*flag.FlagSet
}

type Command interface {
	interfaces.CommandComponent
	Run(Dep)
}

var commands = map[string]Command{}

func Commands() map[string]Command {
	return commands
}

func Register(name string, cmd Command) {
	if _, ok := commands[name]; ok {
		panic("command added more than once: " + name)
	}

	commands[name] = cmd
}
