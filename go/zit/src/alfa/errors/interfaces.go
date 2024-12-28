package errors

type Flusher interface {
	Flush() error
}

type Func func() error

type FuncWithStackInfo struct {
	Func
	StackInfo
}

type WaitGroup interface {
	Do(Func) bool
	DoAfter(Func)
	GetError() error
}
