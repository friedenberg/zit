package commands

import (
	"flag"
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/debug"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type Dependencies struct {
	errors.Context
	config_mutable_cli.Config
}

type Command2 interface {
	GetFlagSet() *flag.FlagSet
	Run(Dependencies) int
}

type CommandWithRepo interface {
	Run(*repo_local.Repo, ...string) error
}

type CommandWithEnv interface {
	Run(*env.Env, ...string)
}

type CommandWithContext interface {
	Run(*repo_local.Repo, ...string)
}

type WithCompletion interface {
	Complete(u *repo_local.Repo, args ...string) (err error)
}

type command struct {
	withoutRepo bool
	Command     CommandWithContext
	*flag.FlagSet
}

func (cmd command) GetFlagSet() *flag.FlagSet {
	return cmd.FlagSet
}

func (cmd command) Run(
	dependencies Dependencies,
) (exitStatus int) {
	// TODO use options when making dirLayout
	var dirLayout dir_layout.Layout

	{
		var err error

		if dirLayout, err = dir_layout.MakeDefault(
			dependencies.Debug,
		); err != nil {
			dependencies.CancelWithError(err)
			return
		}
	}

	// TODO move to env
	if _, err := debug.MakeContext(
		dependencies.Context,
		dependencies.Debug,
	); err != nil {
		dependencies.CancelWithError(err)
		return
	}

	env := env.Make(
		dependencies.Context,
		cmd.GetFlagSet(),
		dependencies.Config,
		dirLayout,
	)

	cmdArgs := cmd.Args()

	var u *repo_local.Repo

	options := repo_local.OptionsEmpty

	if og, ok := cmd.Command.(repo_local.OptionsGetter); ok {
		options = og.GetEnvironmentInitializeOptions()
	}

	{
		var err error

		if u, err = repo_local.Make(
			env,
			options,
		); err != nil {
			dependencies.CancelWithError(err)
			return
		}

		defer errors.DeferredFlusher(&err, u)
	}

	defer func() {
		if err := u.GetRepoLayout().ResetTempOnExit(
			dependencies.Context,
		); err != nil {
			dependencies.CancelWithError(err)
			return
		}
	}()

	switch {
	case u.GetConfig().Complete:
		var t WithCompletion
		haystack := any(cmd.Command)

	LOOP:
		for {
			switch c := haystack.(type) {
			case commandWithResult:
				haystack = c.CommandWithRepo
				continue LOOP

			case WithCompletion:
				t = c
				break LOOP

			default:
				dependencies.Cancel(errors.BadRequestf("Command does not support completion"))
				return
			}
		}

		if err := t.Complete(u, cmdArgs...); err != nil {
			dependencies.CancelWithError(err)
			return
		}

	default:

		func() {
			defer func() {
				// if r := recover(); r != nil {
				// 	result = ErrorResult{error: errors.Errorf("panicked: %s", r)}
				// }
			}()

			cmd.Command.Run(u, cmdArgs...)
		}()
	}

	return
}

var commands = map[string]Command2{}

func Commands() map[string]Command2 {
	return commands
}

func _registerCommand(
	withoutRepo bool,
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
			Command: commandWithResult{CommandWithRepo: mft(f)},
			FlagSet: f,
		}

	case func(*flag.FlagSet) CommandWithContext:
		commands[n] = command{
			withoutRepo: withoutRepo,
			Command:     mft(f),
			FlagSet:     f,
		}

	default:
		panic(fmt.Sprintf("command make func not supported: %T", mft))
	}
}

func registerCommand(n string, makeFunc any) {
	_registerCommand(false, n, makeFunc)
}

func registerCommandWithoutRepo(
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
