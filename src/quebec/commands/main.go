package commands

import (
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/debug"
	"github.com/friedenberg/zit/src/echo/standort"
	"github.com/friedenberg/zit/src/india/erworben"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

func Run(args []string) (exitStatus int) {
	var err error

	defer func() {
		errors.Log().Print("checking for open files")
		l := files.Len()
		errors.Log().Printf("open files: %d", l)

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
		return cmd.PrintUsage(errors.Normalf("No subcommand '%s'", specifiedSubcommand))
	}

	args = os.Args[2:]

	konfigCli := erworben.DefaultCli()
	konfigCli.AddToFlags(cmd.FlagSet)

	if err = cmd.FlagSet.Parse(args); err != nil {
		err = errors.Wrap(err)
		return
	}

	var dc *debug.Context

	if dc, err = debug.MakeContext(konfigCli.Debug); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, dc.Close)

	cmdArgs := cmd.FlagSet.Args()

	var u *umwelt.Umwelt

	if u, err = umwelt.Make(konfigCli); err != nil {
		//the store doesn't exist yet
		switch {
		case errors.IsNotExist(err):
			err = nil

		case errors.Is(err, standort.ErrNotInZitDir{}) && cmd.FlagSet.Name() == "init":
			if err = cmd.Command.Run(u, cmdArgs...); err != nil {
				err = errors.Wrap(err)
				return
			}

			return

		default:
			err = errors.Wrap(err)
			return
		}
	}

	defer errors.Deferred(&err, u.Flush)

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
		if err = cmd.Command.Run(u, cmdArgs...); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
