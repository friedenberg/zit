package main

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/quebec/commands"
)

func main() {
	var exitStatus int

	for {
		ctx := errors.MakeContextDefault()
		ctx.SetCancelOnSIGINT()

		if err := ctx.Run(
			func(ctx *errors.Context) {
				commands.Run(ctx, os.Args...)
			},
		); err != nil {
			var retryable errors.Retryable

			if errors.As(err, &retryable) {
				ui.Err().Print("retryable")
				// TODO retry
				// continue
			}

			var normalError errors.StackTracer
			exitStatus = 1

			if errors.As(err, &normalError) && !normalError.ShouldShowStackTrace() {
				ui.Err().Printf("%s", normalError.Error())
			} else {
				ui.Err().Print(err)
			}
		}

		break
	}

	os.Exit(exitStatus)
}
