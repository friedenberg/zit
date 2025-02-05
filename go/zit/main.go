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
			func(ctx errors.Context) {
				commands.Run(ctx, os.Args...)
			},
		); err != nil {
			var helpful errors.Helpful

			if errors.As(err, &helpful) {
				errors.PrintHelpful(ui.Err(), helpful)
				break
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
