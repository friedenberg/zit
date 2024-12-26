package commands

import (
	"flag"
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type Dependencies struct {
	errors.Context
	config_mutable_cli.Config
}

type CommandWithDependencies interface {
	GetFlagSet() *flag.FlagSet
	RunWithDependencies(Dependencies) int
}

type CommandWithEnv interface {
	RunWithEnv(*env.Env, ...string)
}

type CommandWithRepo interface {
	RunWithRepo(*repo_local.Repo, ...string)
}

type CommandWithQuery interface {
	RunWithQuery(store *repo_local.Repo, ids *query.Group) error
}

type CommandCompletionWithRepo interface {
	CompleteWithRepo(u *repo_local.Repo, args ...string)
}

var commands = map[string]CommandWithDependencies{}

func Commands() map[string]CommandWithDependencies {
	return commands
}

func registerCommand(
	n string,
	makeFunc any,
) {
	f := flag.NewFlagSet(n, flag.ExitOnError)

	if _, ok := commands[n]; ok {
		panic("command added more than once: " + n)
	}

	switch mft := makeFunc.(type) {
	case func(*flag.FlagSet) CommandWithEnv:
		commands[n] = commandWithEnv{
			Command: mft(f),
			FlagSet: f,
		}

	case func(*flag.FlagSet) CommandWithRepo:
		commands[n] = commandWithRepo{
			Command: mft(f),
			FlagSet: f,
		}

	default:
		panic(fmt.Sprintf("command make func not supported: %T", mft))
	}
}

func registerCommandWithoutRepo(
	n string,
	makeFunc any,
) {
	registerCommand(n, makeFunc)
}

func registerCommandWithQuery(
	n string,
	makeFunc func(*flag.FlagSet) CommandWithQuery,
) {
	registerCommand(
		n,
		func(f *flag.FlagSet) CommandWithRepo {
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
	registerCommand(
		n,
		func(f *flag.FlagSet) CommandWithRepo {
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
	registerCommand(
		n,
		func(f *flag.FlagSet) CommandWithRepo {
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
