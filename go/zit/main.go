package main

import (
	"context"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/quebec/commands"
)

func main() {
	ctx := errors.MakeContext(context.Background())
	defer ctx.Cancel(nil)

	ctx.SetCancelOnSIGINT()

	var exitStatus int

	func() {
		defer func() {
			if r := recover(); r != nil {
				exitStatus = 1

				if r != errors.ErrContextCancelled {
					panic(r)
				}
			}
		}()

		exitStatus = commands.Run(ctx, os.Args...)
	}()

	if err := context.Cause(ctx); err != nil {
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
