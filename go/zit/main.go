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
			func(ctx errors.IContext) {
				commands.Run(ctx, os.Args...)
			},
		); err != nil {
			var retryable errors.Retryable

			if errors.As(err, &retryable) {
				// TODO retry
				// continue
			}

			var helpful errors.Helpful

			if errors.As(err, &helpful) {
				ui.Err().Printf("Error: %s", helpful.Error())
				ui.Err().Printf("\nCause:")

				for _, causeLine := range helpful.ErrorCause() {
					ui.Err().Print(causeLine)
				}

				ui.Err().Printf("\nRecovery:")

				for _, recoveryLine := range helpful.ErrorRecovery() {
					ui.Err().Print(recoveryLine)
				}

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
