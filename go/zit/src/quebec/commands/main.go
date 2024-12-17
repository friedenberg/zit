package commands

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/debug"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/golf/mutable_config_blobs"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
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

	cliConfig := mutable_config_blobs.DefaultCli()
	cliConfig.AddToFlags(cmd.FlagSet)

	if err := cmd.Parse(args); err != nil {
		ctx.Cancel(err)
		return
	}

	var primitiveFSHome dir_layout.Primitive
	var err error

	if primitiveFSHome, err = dir_layout.MakePrimitive(cliConfig.Debug); err != nil {
		ctx.Cancel(errors.Wrap(err))
		return
	}

	if _, err = debug.MakeContext(ctx, cliConfig.Debug); err != nil {
		ctx.Cancel(errors.Wrap(err))
		return
	}

	cmdArgs := cmd.Args()

	var u *env.Local

	options := env.OptionsEmpty

	if og, ok := cmd.Command.(env.OptionsGetter); ok {
		options = og.GetEnvironmentInitializeOptions()
	}

	if u, err = env.MakeLocal(
		ctx,
		cmd.FlagSet,
		cliConfig,
		options,
		primitiveFSHome,
	); err != nil {
		if cmd.withoutEnv {
			err = nil
		} else {
			ctx.Cancel(errors.Wrap(err))
			return
		}
	}

	defer u.PrintMatchedArchiviertIfNecessary()
	defer errors.DeferredFlusher(&err, u)

	defer func() {
		if err = u.GetDirectoryLayout().ResetTempOnExit(ctx); err != nil {
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
