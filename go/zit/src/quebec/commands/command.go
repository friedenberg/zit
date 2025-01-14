package commands

import (
	"flag"
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
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

type Command2 interface {
	Run(Dependencies)
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

type CommandWithEnv interface {
	RunWithEnv(*env.Env, ...string)
}

type CommandWithArchive interface {
	RunWithArchive(repo.Archive, ...string)
}

type CommandWithWorkingCopy interface {
	RunWithWorkingCopy(repo.WorkingCopy, ...string)
}

type CommandWithLocalWorkingCopy interface {
	RunWithLocalWorkingCopy(*local_working_copy.Repo, ...string)
}

type CommandWithBlobStore interface {
	RunWithBlobStore(command_components.BlobStoreWithEnv, ...string)
}

type CommandWithQuery interface {
	RunWithQuery(store *local_working_copy.Repo, ids *query.Group)
}

type CommandWithQueryAndBuilderOptions interface {
	query.BuilderOptionGetter
	CommandWithQuery
}

type CommandCompletionWithRepo interface {
	CompleteWithRepo(u *local_working_copy.Repo, args ...string)
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

	case func(*flag.FlagSet) CommandWithEnv:
		commands[n] = commandWithEnv{
			Command: cmd(f),
			FlagSet: f,
		}

	case func(*flag.FlagSet) CommandWithLocalWorkingCopy:
		commands[n] = commandWithLocalWorkingCopy{
			Command: cmd(f),
			FlagSet: f,
		}

	case func(*flag.FlagSet) CommandWithBlobStore:
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
	makeFunc func(*flag.FlagSet) CommandWithQuery,
) {
	registerCommand(
		n,
		func(f *flag.FlagSet) CommandWithLocalWorkingCopy {
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
		func(f *flag.FlagSet) CommandWithLocalWorkingCopy {
			cmd := &commandWithRemoteAndQuery{
				CommandWithRemoteAndQuery: cwraq,
			}

			cmd.SetFlagSet(f)

			return cmd
		},
	)
}
