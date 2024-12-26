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
	cancel context.CancelCauseFunc
}

func MakeContextDefault() Context {
	return MakeContext(context.Background())
}

func MakeContext(in context.Context) Context {
	ctx, cancel := context.WithCancelCause(in)

	return Context{
		Context: ctx,
		cancel:  cancel,
	}
}

func (c Context) Cause() error {
	if err := context.Cause(c.Context); err != nil {
		if Is(err, ErrContextCancelled) {
			return nil
		} else {
			return err
		}
	}

	return nil
}

func (c Context) ContinueOrPanicOnDone() {
	select {
	default:
	case <-c.Context.Done():
		panic(ErrContextCancelled)
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
		c.cancel(Signal{Signal: <-ch})
	}()
}

func (c Context) Run(f func(Context)) error {
	func() {
		defer c.cancel(ErrContextCancelled)
		defer func() {
			if r := recover(); r != nil {
				if r != ErrContextCancelled {
					panic(r)
				}
			}
		}()

		f(c)
	}()

	return c.Cause()
}

// Must executes a function even if the context has been cancelled. If the
// function returns an error, Must cancels the context and offers a heartbeat to
// panic. It is meant for defers that must be executed, like closing files,
// flushing buffers, releasing locks.
func (c Context) Must(f func() error) {
	defer c.ContinueOrPanicOnDone()

	if err := f(); err != nil {
		c.cancel(WrapN(1, err))
	}
}

func (c Context) MustClose(closer io.Closer) {
	c.Must(closer.Close)
}

func (c Context) MustFlush(flusher Flusher) {
	c.Must(flusher.Flush)
}

func (c Context) Cancel() {
	defer c.ContinueOrPanicOnDone()
	c.cancel(ErrContextCancelled)
}

func (c Context) CancelWithError(err error) {
	defer c.ContinueOrPanicOnDone()
	c.cancel(WrapN(1, err))
}

func (c Context) CancelWithErrorAndFormat(err error, f string, values ...any) {
	defer c.ContinueOrPanicOnDone()
	c.cancel(WrapN(1, err))
	c.cancel(
		&stackWrapError{
			StackInfo: MustStackInfo(1),
			error:     fmt.Errorf(f, values...),
			next:      WrapSkip(1, err),
		},
	)
}

func (c Context) CancelWithErrorf(f string, values ...any) {
	defer c.ContinueOrPanicOnDone()
	c.cancel(WrapSkip(1, fmt.Errorf(f, values...)))
}

func (c Context) CancelWithBadRequestf(f string, values ...any) {
	defer c.ContinueOrPanicOnDone()
	c.cancel(&errBadRequest{xerrors.Errorf(f, values...)})
}
