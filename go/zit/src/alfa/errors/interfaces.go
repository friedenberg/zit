package errors

type Flusher interface {
	Flush() error
}

type Func func() error

type ErrorWaitGroup interface {
	Do(Func) bool
	DoAfter(Func)
	GetError() error
}
