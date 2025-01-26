package errors

import (
	ConTeXT "context"
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

type Context interface {
	ConTeXT.Context

	Cause() error
	Continue() bool
	ContinueOrPanicOnDone()
	SetCancelOnSIGINT()
	SetCancelOnSIGHUP()
	SetCancelOnSignals(signals ...os.Signal)
	Run(f func(Context)) error

	// `After` runs a function after the context is complete (regardless of any
	// errors). `After`s are run in the reverse order of when they are called, like
	// defers but on a whole-program level.
	After(f func() error)
	AfterWithContext(f func(Context) error)

	// `Must` executes a function even if the context has been cancelled. If the
	// function returns an error, `Must` cancels the context and offers a heartbeat to
	// panic. It is meant for defers that must be executed, like closing files,
	// flushing buffers, releasing locks.
	Must(f func() error)
	MustWithContext(f func(Context) error)
	MustClose(closer io.Closer)
	MustFlush(flusher Flusher)
	Cancel()

	CancelWithError(err error)
	CancelWithErrorAndFormat(err error, f string, values ...any)
	CancelWithErrorf(f string, values ...any)
	CancelWithBadRequestf(f string, values ...any)
	CancelWithNotImplemented()
}

type context struct {
	ConTeXT.Context
	cancelFunc ConTeXT.CancelCauseFunc

	signals chan os.Signal

	lock          sync.Mutex
	doAfter       []FuncWithStackInfo
	doAfterErrors []error // TODO expose and use
}

func MakeContextDefault() *context {
	return MakeContext(ConTeXT.Background())
}

func MakeContext(in ConTeXT.Context) *context {
	ctx, cancel := ConTeXT.WithCancelCause(in)

	return &context{
		Context:    ctx,
		cancelFunc: cancel,
		signals:    make(chan os.Signal, 1),
	}
}

func (c *context) Cause() error {
	if err := ConTeXT.Cause(c.Context); err != nil {
		if Is(err, errContextCancelled) {
			return nil
		} else {
			return err
		}
	}

	return nil
}

func (c *context) Continue() bool {
	select {
	default:
		return true

	case <-c.Done():
		return false
	}
}

func (c *context) ContinueOrPanicOnDone() {
	if !c.Continue() {
		panic(errContextCancelled)
	}
}

func (c *context) SetCancelOnSIGINT() {
	c.SetCancelOnSignals(syscall.SIGINT)
}

func (c *context) SetCancelOnSIGHUP() {
	c.SetCancelOnSignals(syscall.SIGHUP)
}

func (c *context) SetCancelOnSignals(signals ...os.Signal) {
	signal.Notify(c.signals, signals...)
}

func (c *context) Run(f func(Context)) error {
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
				if err, ok := r.(error); !ok {
					panic(r)
				} else {
					c.cancel(err)
				}
				// fmt.Printf("%s", debug.Stack())
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

func (c *context) cancel(err error) {
	c.cancelFunc(err)
}

//go:noinline
func (c *context) after(skip int, f func() error) {
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
func (c *context) After(f func() error) {
	c.after(1, f)
}

//go:noinline
func (c *context) AfterWithContext(f func(Context) error) {
	c.after(1, func() error { return f(c) })
}

// `Must` executes a function even if the context has been cancelled. If the
// function returns an error, `Must` cancels the context and offers a heartbeat to
// panic. It is meant for defers that must be executed, like closing files,
// flushing buffers, releasing locks.
func (c *context) Must(f func() error) {
	defer c.ContinueOrPanicOnDone()

	if err := f(); err != nil {
		c.cancel(WrapN(1, err))
	}
}

func (c *context) MustWithContext(f func(Context) error) {
	defer c.ContinueOrPanicOnDone()

	if err := f(c); err != nil {
		c.cancel(WrapN(1, err))
	}
}

func (c *context) MustClose(closer io.Closer) {
	c.Must(closer.Close)
}

func (c *context) MustFlush(flusher Flusher) {
	c.Must(flusher.Flush)
}

// TODO make this private and part of the run method
func (c *context) Cancel() {
	defer c.ContinueOrPanicOnDone()
	c.cancelWithoutPanic()
}

func (c *context) cancelWithoutPanic() {
	c.cancel(errContextCancelled)
}

func (c *context) CancelWithError(err error) {
	defer c.ContinueOrPanicOnDone()
	c.cancel(WrapN(1, err))
}

func (c *context) CancelWithErrorAndFormat(err error, f string, values ...any) {
	defer c.ContinueOrPanicOnDone()
	c.cancel(
		&stackWrapError{
			StackInfo: MustStackInfo(1),
			error:     fmt.Errorf(f, values...),
			next:      WrapSkip(1, err),
		},
	)
}

func (c *context) CancelWithErrorf(f string, values ...any) {
	defer c.ContinueOrPanicOnDone()
	c.cancel(WrapSkip(1, fmt.Errorf(f, values...)))
}

func (c *context) CancelWithBadRequestf(f string, values ...any) {
	defer c.ContinueOrPanicOnDone()
	c.cancel(&errBadRequest{xerrors.Errorf(f, values...)})
}

func (c *context) CancelWithNotImplemented() {
	defer c.ContinueOrPanicOnDone()
	c.cancel(ErrNotImplemented)
}
