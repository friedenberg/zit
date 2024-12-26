package errors

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/xerrors"
)

var ErrContextCancelled = New("context cancelled")

type Context struct {
	context.Context
	Cancel context.CancelCauseFunc
}

func MakeContextDefault() Context {
	return MakeContext(context.Background())
}

func MakeContext(in context.Context) Context {
	ctx, cancel := context.WithCancelCause(in)

	return Context{
		Context: ctx,
		Cancel:  cancel,
	}
}

func (c Context) Cause() error {
	if err := context.Cause(c); err != nil {
		if Is(err, ErrContextCancelled) {
			return nil
		} else {
			return err
		}
	}

	return nil
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
		c.Cancel(WrapN(1, err))
	}

	c.Heartbeat()
}

func (c Context) Closer(closer io.Closer) {
	c.Must(closer.Close)
}

func (c Context) Flusher(flusher Flusher) {
	c.Must(flusher.Flush)
}

func (c Context) CancelWithError(err error) {
	c.Cancel(WrapN(1, err))
	panic(ErrContextCancelled)
}

func (c Context) CancelWithErrorf(f string, values ...any) {
	c.Cancel(WrapSkip(1, fmt.Errorf(f, values...)))
	panic(ErrContextCancelled)
}

func (c Context) CancelWithBadRequestf(f string, values ...any) {
	c.Cancel(&errBadRequest{xerrors.Errorf(f, values...)})
	panic(ErrContextCancelled)
}
