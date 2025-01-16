package errors

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/xerrors"
)

var errContextCancelled errContextCancelledExpected

type errContextCancelledExpected struct {
	error
}

func (err errContextCancelledExpected) Error() string {
	if err.error == nil {
		return "context cancelled"
	} else {
		return fmt.Sprintf("context cancelled: %s", err.error)
	}
}

func (err errContextCancelledExpected) Is(target error) bool {
	_, ok := target.(errContextCancelledExpected)
	return ok
}

type Context struct {
	context.Context
	cancelFunc context.CancelCauseFunc

	signals chan os.Signal

	lock          sync.Mutex
	doAfter       []FuncWithStackInfo
	doAfterErrors []error // TODO expose and use
}

func MakeContextDefault() *Context {
	return MakeContext(context.Background())
}

func MakeContext(in context.Context) *Context {
	ctx, cancel := context.WithCancelCause(in)

	return &Context{
		Context:    ctx,
		cancelFunc: cancel,
		signals:    make(chan os.Signal, 1),
	}
}

func (c *Context) Cause() error {
	if err := context.Cause(c.Context); err != nil {
		if Is(err, errContextCancelled) {
			return nil
		} else {
			return err
		}
	}

	return nil
}

func (c *Context) Continue() bool {
	select {
	default:
		return true

	case <-c.Done():
		return false
	}
}

func (c *Context) ContinueOrPanicOnDone() {
	if !c.Continue() {
		panic(errContextCancelled)
	}
}

func (c *Context) SetCancelOnSIGINT() {
	c.SetCancelOnSignals(syscall.SIGINT)
}

func (c *Context) SetCancelOnSIGHUP() {
	c.SetCancelOnSignals(syscall.SIGHUP)
}

func (c *Context) SetCancelOnSignals(signals ...os.Signal) {
	signal.Notify(c.signals, signals...)
}

func (c *Context) Run(f func(*Context)) error {
	go func() {
		select {
		case <-c.Done():
		case sig := <-c.signals:
			c.cancel(errContextCancelledExpected{Signal{Signal: sig}})
		}

		signal.Stop(c.signals)
	}()

	func() {
		defer c.cancel(errContextCancelled)
		defer func() {
			if r := recover(); r != nil {
				var err error

				{
					var ok bool

					if err, ok = r.(error); !ok {
						panic(r)
					}
				}

				if !Is(err, errContextCancelledExpected{}) {
					panic(err)
				}
			}
		}()

		f(c)
	}()

	for i := len(c.doAfter) - 1; i >= 0; i-- {
		doAfter := c.doAfter[i]
		err := doAfter.Func()
		if err != nil {
			c.doAfterErrors = append(
				c.doAfterErrors,
				doAfter.Wrap(err),
			)
		}
	}

	return c.Cause()
}

func (c *Context) cancel(err error) {
	c.cancelFunc(err)
}

//go:noinline
func (c *Context) after(skip int, f func() error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	si, _ := MakeStackInfo(skip + 1)

	c.doAfter = append(
		c.doAfter,
		FuncWithStackInfo{
			Func:      f,
			StackInfo: si,
		},
	)
}

// `After` runs a function after the context is complete (regardless of any
// errors). `After`s are run in the reverse order of when they are called, like
// defers but on a whole-program level.
//
//go:noinline
func (c *Context) After(f func() error) {
	c.after(1, f)
}

//go:noinline
func (c *Context) AfterWithContext(f func(*Context) error) {
	c.after(1, func() error { return f(c) })
}

// `Must` executes a function even if the context has been cancelled. If the
// function returns an error, `Must` cancels the context and offers a heartbeat to
// panic. It is meant for defers that must be executed, like closing files,
// flushing buffers, releasing locks.
func (c *Context) Must(f func() error) {
	defer c.ContinueOrPanicOnDone()

	if err := f(); err != nil {
		c.cancel(WrapN(1, err))
	}
}

func (c *Context) MustWithContext(f func(*Context) error) {
	defer c.ContinueOrPanicOnDone()

	if err := f(c); err != nil {
		c.cancel(WrapN(1, err))
	}
}

func (c *Context) MustClose(closer io.Closer) {
	c.Must(closer.Close)
}

func (c *Context) MustFlush(flusher Flusher) {
	c.Must(flusher.Flush)
}

// TODO make this private and part of the run method
func (c *Context) Cancel() {
	defer c.ContinueOrPanicOnDone()
	c.cancelWithoutPanic()
}

func (c *Context) cancelWithoutPanic() {
	c.cancel(errContextCancelled)
}

func (c *Context) CancelWithError(err error) {
	defer c.ContinueOrPanicOnDone()
	c.cancel(WrapN(1, err))
}

func (c *Context) CancelWithErrorAndFormat(err error, f string, values ...any) {
	defer c.ContinueOrPanicOnDone()
	c.cancel(
		&stackWrapError{
			StackInfo: MustStackInfo(1),
			error:     fmt.Errorf(f, values...),
			next:      WrapSkip(1, err),
		},
	)
}

func (c *Context) CancelWithErrorf(f string, values ...any) {
	defer c.ContinueOrPanicOnDone()
	c.cancel(WrapSkip(1, fmt.Errorf(f, values...)))
}

func (c *Context) CancelWithBadRequestf(f string, values ...any) {
	defer c.ContinueOrPanicOnDone()
	c.cancel(&errBadRequest{xerrors.Errorf(f, values...)})
}

func (c *Context) CancelWithNotImplemented() {
	defer c.ContinueOrPanicOnDone()
	c.cancel(ErrNotImplemented)
}
