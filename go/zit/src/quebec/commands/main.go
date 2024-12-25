package commands

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/debug"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

// TODO switch to returning result
func Run(
	ctx errors.Context,
	args ...string,
) (exitStatus int) {
	var cmd command

	if len(os.Args) < 1 {
		ui.Log().Print("printing usage")
		return cmd.PrintUsage(nil)
	}

	if len(os.Args) == 1 {
		return cmd.PrintUsage(errors.Errorf("No subcommand provided."))
	}

	cmds := Commands()
	specifiedSubcommand := os.Args[1]

	ok := false

	if cmd, ok = cmds[specifiedSubcommand]; !ok {
		return cmd.PrintUsage(
			errors.BadRequestf("No subcommand '%s'", specifiedSubcommand),
		)
	}

	args = os.Args[2:]

	configCli := config_mutable_cli.Default()
	configCli.AddToFlags(cmd.FlagSet)

	if err := cmd.Parse(args); err != nil {
		ctx.Cancel(err)
		return
	}

	var dirLayout dir_layout.Layout
	var err error

	if dirLayout, err = dir_layout.MakeDefault(
		configCli.Debug,
	); err != nil {
		ctx.Cancel(errors.Wrap(err))
		return
	}

	if _, err = debug.MakeContext(ctx, configCli.Debug); err != nil {
		ctx.Cancel(errors.Wrap(err))
		return
	}

	cmdArgs := cmd.Args()

	var u *repo_local.Repo

	options := repo_local.OptionsEmpty

	if og, ok := cmd.Command.(repo_local.OptionsGetter); ok {
		options = og.GetEnvironmentInitializeOptions()
	}

	env := env.Make(
		ctx,
		cmd.FlagSet,
		configCli,
		dirLayout,
	)

	if u, err = repo_local.Make(
		env,
		options,
	); err != nil {
		if cmd.withoutRepo {
			err = nil
		} else {
			ctx.Cancel(errors.Wrap(err))
			return
		}
	}

	defer errors.DeferredFlusher(&err, u)

	defer func() {
		if err = u.GetRepoLayout().ResetTempOnExit(ctx); err != nil {
			ctx.Cancel(errors.Wrap(err))
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
				haystack = c.Command
				continue LOOP

			case WithCompletion:
				t = c
				break LOOP

			default:
				ctx.Cancel(errors.BadRequestf("Command does not support completion"))
				return
			}
		}

		if err := t.Complete(u, cmdArgs...); err != nil {
			ctx.Cancel(err)
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
