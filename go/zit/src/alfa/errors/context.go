package errors

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

type SignalError struct {
	os.Signal
}

func (err SignalError) Error() string {
	return fmt.Sprintf("signal: %q", err.Signal)
}

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

func (c Context) Heartbeat() {
	select {
	default:
		return

	case <-c.Done():
		if err := context.Cause(c); err != nil {
			panic(err)
		}
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
		c.Cancel(SignalError{Signal: <-ch})
	}()
}

func (c Context) Closer(
	closer io.Closer,
) {
	if err := closer.Close(); err != nil {
		c.Cancel(err)
	}
}
