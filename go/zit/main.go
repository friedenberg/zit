package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/quebec/commands"
)

func main() {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	makeSIGINTWatchChannelAndCancelContextIfNecessary(cancel)

	exitStatus := commands.Run(ctx, cancel, os.Args...)

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

func makeSIGINTWatchChannelAndCancelContextIfNecessary(
	cancel context.CancelCauseFunc,
) {
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGINT)

	go func() {
		<-ch
		cancel(nil)
		os.Exit(1)
	}()
}
