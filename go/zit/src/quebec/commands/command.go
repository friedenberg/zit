package commands

import (
	"flag"
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

type Dependencies struct {
	*errors.Context
	config_mutable_cli.Config
}

type Command interface {
	GetFlagSet() *flag.FlagSet
	Run(Dependencies)
}

type CommandWithEnv interface {
	RunWithEnv(*env.Env, ...string)
}

type CommandWithRepo interface {
	RunWithRepo(*repo_local_working_copy.Repo, ...string)
}

type CommandWithBlobStore interface {
	RunWithBlobStore(command_components.BlobStoreWithEnv, ...string)
}

type CommandWithQuery interface {
	RunWithQuery(store *repo_local_working_copy.Repo, ids *query.Group)
}

type CommandWithQueryAndBuilderOptions interface {
	query.BuilderOptionGetter
	CommandWithQuery
}

type CommandCompletionWithRepo interface {
	CompleteWithRepo(u *repo_local_working_copy.Repo, args ...string)
}

var commands = map[string]Command{}

func Commands() map[string]Command {
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
	case func(*flag.FlagSet) Command:
		commands[n] = mft(f)

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

	case func(*flag.FlagSet) CommandWithBlobStore:
		commands[n] = commandWithBlobStore{
			Command: mft(f),
			FlagSet: f,
		}

	default:
		panic(fmt.Sprintf("command make func not supported: %T", mft))
	}
}

func registerCommandWithQuery(
	n string,
	makeFunc func(*flag.FlagSet) CommandWithQuery,
) {
	registerCommand(
		n,
		func(f *flag.FlagSet) CommandWithRepo {
			cmd := &commandWithQuery{
				CommandWithQuery: makeFunc(f),
			}

			cmd.SetFlagSet(f)

			return cmd
		},
	)
}

func registerCommandWithRemoteAndQuery(
	n string,
	cwraq CommandWithRemoteAndQuery,
) {
	registerCommand(
		n,
		func(f *flag.FlagSet) CommandWithRepo {
			cmd := &commandWithRemoteAndQuery{
				CommandWithRemoteAndQuery: cwraq,
			}

			cmd.SetFlagSet(f)

			return cmd
		},
	)
}

// TODO explore a different flag and command model
func registerCommandWithFlags(
	n string,
	cmd interfaces.CommandComponent,
) {
	registerCommand(
		n,
		func(f *flag.FlagSet) interfaces.CommandComponent {
			cmd.SetFlagSet(f)
			return cmd
		},
	)
}
