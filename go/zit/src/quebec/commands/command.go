package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type Command interface {
	Run(*env.Env, ...string) error
}

type WithCompletion interface {
	Complete(u *env.Env, args ...string) (err error)
}

type command struct {
	withoutEnv bool
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

func registerCommandWithoutEnvironment(
	n string,
	makeFunc func(*flag.FlagSet) Command,
) {
	f := flag.NewFlagSet(n, flag.ExitOnError)

	c := makeFunc(f)

	if _, ok := commands[n]; ok {
		panic("command added more than once: " + n)
	}

	commands[n] = command{
		withoutEnv: true,
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

	cweq := &commandWithQuery{
		CommandWithQuery: c,
	}

	f.Var(&cweq.RepoId, "kasten", "none or Browser")
	f.BoolVar(&cweq.ExcludeUntracked, "exclude-untracked", false, "")
	f.BoolVar(&cweq.ExcludeRecognized, "exclude-recognized", false, "")

	co := command{
		Command: cweq,
		FlagSet: f,
	}

	commands[n] = co
}
