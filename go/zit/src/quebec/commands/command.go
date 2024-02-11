package commands

import (
	"flag"

	"code.linenisgreat.com/zit/src/oscar/umwelt"
)

type Command interface {
	Run(*umwelt.Umwelt, ...string) error
}

type WithCompletion interface {
	Complete(u *umwelt.Umwelt, args ...string) (err error)
}

type command struct {
	sansUmwelt bool
	Command
	*flag.FlagSet
}

var commands = map[string]command{}

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

func registerCommandSansUmwelt(n string, makeFunc func(*flag.FlagSet) Command) {
	f := flag.NewFlagSet(n, flag.ExitOnError)

	c := makeFunc(f)

	if _, ok := commands[n]; ok {
		panic("command added more than once: " + n)
	}

	commands[n] = command{
		sansUmwelt: true,
		Command:    c,
		FlagSet:    f,
	}

	return
}

func registerCommandWithCwdQuery(
	n string,
	makeFunc func(*flag.FlagSet) CommandWithCwdQuery,
) {
	f := flag.NewFlagSet(n, flag.ExitOnError)

	c := makeFunc(f)

	if _, ok := commands[n]; ok {
		panic("command added more than once: " + n)
	}

	commands[n] = command{
		Command: commandWithCwdQuery{CommandWithCwdQuery: c},
		FlagSet: f,
	}

	return
}

func registerCommandWithQuery(
	n string,
	makeFunc func(*flag.FlagSet) CommandWithQuery,
) {
	f := flag.NewFlagSet(n, flag.ExitOnError)

	c := makeFunc(f)

	if _, ok := commands[n]; ok {
		panic("command added more than once: " + n)
	}

	commands[n] = command{
		Command: commandWithQuery{
			CommandWithQuery: c,
		},
		FlagSet: f,
	}

	return
}
