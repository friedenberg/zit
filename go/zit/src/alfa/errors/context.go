package errors

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type Context struct {
	context.Context
	Cancel context.CancelCauseFunc
}

func MakeContext(in context.Context) Context {
	ctx, cancel := context.WithCancelCause(in)

	return Context{
		Context: ctx,
		Cancel:  cancel,
	}
}

func (c Context) SetCancelOnSIGINT() {
	c.SetCancelOnSignals(syscall.SIGINT)
}

func (c Context) SetCancelOnSignals(
	signals ...os.Signal,
) {
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, signals...)

	go func() {
		<-ch
		c.Cancel(nil)
	}()
}
