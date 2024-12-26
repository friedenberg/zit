package main

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/quebec/commands"
)

func main() {
	ctx := errors.MakeContextDefault()
	ctx.SetCancelOnSIGINT()

	var exitStatus int

	func() {
		defer ctx.Cancel(errors.ErrContextCancelled)
		defer func() {
			if r := recover(); r != nil {
				exitStatus = 1

				if r != errors.ErrContextCancelled {
					panic(r)
				}
			}
		}()

		commands.Run(ctx, os.Args...)
	}()

	if err := ctx.Cause(); err != nil {
		var normalError errors.StackTracer
		exitStatus = 1

		if errors.As(err, &normalError) && !normalError.ShouldShowStackTrace() {
			ui.Err().Printf("%s", normalError.Error())
		} else {
			ui.Err().Print(err)
		}
	}

	os.Exit(exitStatus)
}
