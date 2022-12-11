package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type Command interface {
	Run(*umwelt.Umwelt, ...string) error
}

type WithCompletion interface {
	Complete(u *umwelt.Umwelt, args ...string) (err error)
}

type command struct {
	Command
	*flag.FlagSet
}

var (
	commands = map[string]command{}
)

func Commands() map[string]command {
	return commands
}

func registerCommand(n string, makeFunc func(*flag.FlagSet) Command) {
	f := flag.NewFlagSet(n, flag.ExitOnError)

	c := makeFunc(f)

	if _, ok := commands[n]; ok {
		panic("command added more than once: " + n)
	}

	commands[n] = command{
		Command: c,
		FlagSet: f,
	}

	return
}
