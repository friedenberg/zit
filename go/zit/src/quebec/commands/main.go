package commands

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/debug"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/mutable_config"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

// TODO switch to returning result
func Run(args []string) (exitStatus int) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGINT)

	go func() {
		<-ch
		cancel()
		os.Exit(1)
	}()

	var err error

	defer func() {
		var normalError errors.StackTracer

		if err != nil {
			exitStatus = 1
		}

		if errors.As(err, &normalError) {
			ui.Err().Printf("%s", normalError.Error())
		} else {
			if err != nil {
				ui.Err().Print(err)
			}
		}
	}()

	var cmd command

	if err != nil {
		err = errors.Wrap(err)
		return
	}

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

	konfigCli := mutable_config.DefaultCli()
	konfigCli.AddToFlags(cmd.FlagSet)

	if err = cmd.Parse(args); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = debug.MakeContext(ctx, konfigCli.Debug); err != nil {
		err = errors.Wrap(err)
		return
	}

	cmdArgs := cmd.Args()

	var u *env.Env

	options := env.OptionsEmpty

	if og, ok := cmd.Command.(env.OptionsGetter); ok {
		options = og.GetEnvironmentInitializeOptions()
	}

	if u, err = env.Make(cmd.FlagSet, konfigCli, options); err != nil {
		if cmd.withoutEnv {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	defer u.PrintMatchedArchiviertIfNecessary()
	defer errors.DeferredFlusher(&err, u)

	var result Result

	defer func() {
		if err = u.GetFSHome().ResetTempOnExit(result.Error); err != nil {
			err = errors.Wrap(err)
			return
		}
	}()

OUTER:
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
				result.Error = errors.BadRequestf("Command does not support completion")
				break OUTER
			}
		}

		result.Error = t.Complete(u, cmdArgs...)

	default:

		func() {
			defer func() {
				// if r := recover(); r != nil {
				// 	result = ErrorResult{error: errors.Errorf("panicked: %s", r)}
				// }
			}()

			result = cmd.Command.Run(u, cmdArgs...)
		}()
	}

	exitStatus = result.ExitCode

	if result.Error != nil {
		exitStatus = 1
		// TODO switch to Err() and update tests
		ui.Out().Print(result.Error)
	}

	return
}
