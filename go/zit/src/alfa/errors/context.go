package errors

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type (
	Context struct {
		context.Context
	}

	ContextOrdinary    Context // "main"
	ContextInterrupted Context // SIGINT
	ContextError       Context // error discovered
	ContextFlushing    Context // coordinating flushing / writing
)

func MakeSIGINTWatchChannelAndCancelContextIfNecessary(
	cancel context.CancelCauseFunc,
) {
	MakeSignalWatchChannelAndCancelContextIfNecessary(cancel, syscall.SIGINT)
}

func MakeSignalWatchChannelAndCancelContextIfNecessary(
	cancel context.CancelCauseFunc,
	signals ...os.Signal,
) {
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, signals...)

	go func() {
		<-ch
		cancel(nil)
		os.Exit(1)
	}()
}

func test(ctx context.Context) {
	chCompletedWithAfterFunc := make(chan struct{})

	completedWithoutAfterFunc := context.AfterFunc(
		ctx,
		func() { // called on <-ctx.Done()
			// TODO do cleanup work
			close(chCompletedWithAfterFunc)
		},
	)

	// TODO do main work (must short-circuit on <-ctx.Done())

	if !completedWithoutAfterFunc() {
		<-chCompletedWithAfterFunc
		// TODO do reset after cleanup or main work
	}
}
