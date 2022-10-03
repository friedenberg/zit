package commands

import (
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/konfig"
	"github.com/friedenberg/zit/src/kilo/umwelt"
)

func Run(args []string) (exitStatus int) {
	var err error

	defer func() {
		errors.Print("checking for open files")
		l := files.Len()
		errors.Printf("open files: %d", l)

		var normalError errors.StackTracer

		if err != nil {
			exitStatus = 1
		}

		if errors.As(err, &normalError) {
			errors.PrintErrf("%s", normalError.Error())
		} else {
			if err != nil {
				errors.PrintErr(err)
			}
		}
	}()

	var cmd command

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if len(os.Args) < 1 {
		errors.Print("printing usage")
		return cmd.PrintUsage(nil)
	}

	if len(os.Args) == 1 {
		return cmd.PrintUsage(errors.Errorf("No subcommand provided."))
	}

	cmds := Commands()
	specifiedSubcommand := os.Args[1]

	ok := false

	if cmd, ok = cmds[specifiedSubcommand]; !ok {
		return cmd.PrintUsage(errors.Errorf("No subcommand '%s'", specifiedSubcommand))
	}

	args = os.Args[2:]

	konfigCli := konfig.DefaultCli()
	konfigCli.AddToFlags(cmd.FlagSet)

	if err = cmd.FlagSet.Parse(args); err != nil {
		err = errors.Wrap(err)
		return
	}

	if konfigCli.Debug {
		df := cmd.SetDebug()
		defer df()
	}

	var k konfig.Konfig

	if k, err = konfigCli.Konfig(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var u *umwelt.Umwelt

	if u, err = umwelt.Make(k); err != nil {
		//the store doesn't exist yet
		if errors.IsNotExist(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	defer u.Flush()

	cmdArgs := cmd.FlagSet.Args()

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
