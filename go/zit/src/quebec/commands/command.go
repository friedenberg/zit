package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
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
}

func registerCommandWithExternalQuery(
	n string,
	makeFunc func(*flag.FlagSet) CommandWithExternalQuery,
) {
	f := flag.NewFlagSet(n, flag.ExitOnError)

	c := makeFunc(f)

	if _, ok := commands[n]; ok {
		panic("command added more than once: " + n)
	}

	cweq := &commandWithExternalQuery{
		CommandWithExternalQuery: c,
	}

	f.Var(&cweq.Kasten, "kasten", "none or Chrome")
	f.BoolVar(&cweq.ExcludeUntracked, "exclude-untracked", false, "")

	co := command{
		Command: cweq,
		FlagSet: f,
	}

	commands[n] = co
}
