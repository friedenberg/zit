package commands

import (
	"os"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/delta/debug"
	"code.linenisgreat.com/zit/src/foxtrot/erworben"
	"code.linenisgreat.com/zit/src/november/umwelt"
)

func Run(args []string) (exitStatus int) {
	var err error

	defer func() {
		var normalError errors.StackTracer

		if err != nil {
			exitStatus = 1
		}

		if errors.As(err, &normalError) {
			errors.Err().Printf("%s", normalError.Error())
		} else {
			if err != nil {
				errors.Err().Print(err)
			}
		}
	}()

	var cmd command

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if len(os.Args) < 1 {
		errors.Log().Print("printing usage")
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
			errors.Normalf("No subcommand '%s'", specifiedSubcommand),
		)
	}

	args = os.Args[2:]

	konfigCli := erworben.DefaultCli()
	konfigCli.AddToFlags(cmd.FlagSet)

	if err = cmd.Parse(args); err != nil {
		err = errors.Wrap(err)
		return
	}

	var dc *debug.Context

	if dc, err = debug.MakeContext(konfigCli.Debug); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, dc)

	cmdArgs := cmd.Args()

	var u *umwelt.Umwelt

	options := umwelt.OptionsEmpty

	if og, ok := cmd.Command.(umwelt.OptionsGetter); ok {
		options = og.GetUmweltInitializeOptions()
	}

	if u, err = umwelt.Make(konfigCli, options); err != nil {
		if cmd.sansUmwelt {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	defer u.PrintMatchedArchiviertIfNecessary()
	defer errors.DeferredFlusher(&err, u)

	switch {
	case u.Konfig().Complete:
		var t WithCompletion
		ok := false

		if t, ok = cmd.Command.(WithCompletion); !ok {
			err = errors.Normalf("Command does not support completion")
			return
		}

		if err = t.Complete(u, cmdArgs...); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		if err = cmd.Run(u, cmdArgs...); err != nil {
			return
		}
	}

	if err == nil {
		if err = u.Standort().ResetTemp(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
