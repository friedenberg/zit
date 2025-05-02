package errors

type Flusher interface {
	Flush() error
}

type (
	Func func() error
)

type FuncWithStackInfo struct {
	Func
	StackFrame
}

type WithStackInfo[T any] struct {
	Contents T
	StackFrame
}

type WaitGroup interface {
	Do(Func) bool
	DoAfter(Func)
	GetError() error
}

func MakeNilFunc(in func()) Func {
	return func() error {
		in()
		return nil
	}
}
