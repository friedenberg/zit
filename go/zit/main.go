package main

import (
	"context"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/quebec/commands"
)

func main() {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	errors.MakeSIGINTWatchChannelAndCancelContextIfNecessary(cancel)

	exitStatus := commands.Run(
		errors.ContextOrdinary{
			Context: ctx,
		},
		cancel,
		os.Args...,
	)

	if err := context.Cause(ctx); err != nil {
		var normalError errors.StackTracer

		if err != nil {
			exitStatus = 1
		}

		if errors.As(err, &normalError) && !normalError.ShouldShowStackTrace() {
			ui.Err().Printf("%s", normalError.Error())
		} else {
			if err != nil {
				ui.Err().Print(err)
			}
		}
	}

	os.Exit(exitStatus)
}
