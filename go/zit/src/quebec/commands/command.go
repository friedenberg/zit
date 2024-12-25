package commands

import (
	"flag"
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type Command interface {
	Run(*repo_local.Local, ...string) error
}

type CommandWithContext interface {
	Run(*repo_local.Local, ...string)
}

type WithCompletion interface {
	Complete(u *repo_local.Local, args ...string) (err error)
}

type command struct {
	withoutEnv bool
	Command    CommandWithContext
	*flag.FlagSet
}

var commands = map[string]command{}

func Commands() map[string]command {
	return commands
}

func _registerCommand(
	env bool,
	n string,
	makeFunc any,
) {
	f := flag.NewFlagSet(n, flag.ExitOnError)

	if _, ok := commands[n]; ok {
		panic("command added more than once: " + n)
	}

	switch mft := makeFunc.(type) {
	case func(*flag.FlagSet) Command:
		commands[n] = command{
			withoutEnv: env,
			Command:    commandWithResult{Command: mft(f)},
			FlagSet:    f,
		}

	case func(*flag.FlagSet) CommandWithContext:
		commands[n] = command{
			withoutEnv: env,
			Command:    mft(f),
			FlagSet:    f,
		}

	default:
		panic(fmt.Sprintf("command make func not supported: %T", mft))
	}
}

func registerCommand(n string, makeFunc any) {
	_registerCommand(false, n, makeFunc)
}

func registerCommandWithoutEnvironment(
	n string,
	makeFunc any,
) {
	_registerCommand(true, n, makeFunc)
}

func registerCommandWithQuery(
	n string,
	makeFunc func(*flag.FlagSet) CommandWithQuery,
) {
	_registerCommand(
		false,
		n,
		func(f *flag.FlagSet) Command {
			cweq := &commandWithQuery{}

			f.Var(&cweq.RepoId, "kasten", "none or Browser")
			f.BoolVar(&cweq.ExcludeUntracked, "exclude-untracked", false, "")
			f.BoolVar(&cweq.ExcludeRecognized, "exclude-recognized", false, "")

			cweq.CommandWithQuery = makeFunc(f)

			return cweq
		},
	)
}

func registerCommandWithRemoteAndQuery(
	n string,
	makeFunc func(*flag.FlagSet) CommandWithRemoteAndQuery,
) {
	_registerCommand(
		false,
		n,
		func(f *flag.FlagSet) CommandWithContext {
			c := &commandWithRemoteAndQuery{}

			f.Var(&c.RepoId, "kasten", "none or Browser")
			f.BoolVar(&c.ExcludeUntracked, "exclude-untracked", false, "")
			f.BoolVar(&c.ExcludeRecognized, "exclude-recognized", false, "")
			f.StringVar(&c.TheirXDGDotenv, "xdg-dotenv", "", "")
			f.BoolVar(&c.UseSocket, "use-socket", false, "")

			c.CommandWithRemoteAndQuery = makeFunc(f)

			return c
		},
	)
}

func registerCommandWithRemoteAndQueryAndWithoutEnvironment(
	n string,
	makeFunc func(*flag.FlagSet) CommandWithRemoteAndQuery,
) {
	_registerCommand(
		true,
		n,
		func(f *flag.FlagSet) CommandWithContext {
			c := &commandWithRemoteAndQuery{}

			f.Var(&c.RepoId, "kasten", "none or Browser")
			f.BoolVar(&c.ExcludeUntracked, "exclude-untracked", false, "")
			f.BoolVar(&c.ExcludeRecognized, "exclude-recognized", false, "")
			f.StringVar(&c.TheirXDGDotenv, "xdg-dotenv", "", "")
			f.BoolVar(&c.UseSocket, "use-socket", false, "")

			c.CommandWithRemoteAndQuery = makeFunc(f)

			return c
		},
	)
}
