package errors

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"
)

var ErrContextCancelled = New("context cancelled")

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
			panic(ErrContextCancelled)
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
		c.Cancel(Signal{Signal: <-ch})
	}()
}

func (c Context) Must(f func() error) {
	if err := f(); err != nil {
		c.Cancel(err)
	}

	c.Heartbeat()
}

func (c Context) Closer(closer io.Closer) {
	c.Must(closer.Close)
}

func (c Context) Flusher(flusher Flusher) {
	c.Must(flusher.Flush)
}
