package errors

type Iser interface {
	error
	Is(error) bool
}

type Unwrapper interface {
	error
	Unwrap() error
}

type Flusher interface {
	Flush() error
}
