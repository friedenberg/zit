package errors

type Flusher interface {
	Flush() error
}

type (
	Func        func() error
	FuncContext func(*Context) error
)

type FuncWithStackInfo struct {
	Func
	StackInfo
}

type WithStackInfo[T any] struct {
	Contents T
	StackInfo
}

type WaitGroup interface {
	Do(Func) bool
	DoAfter(Func)
	GetError() error
}
