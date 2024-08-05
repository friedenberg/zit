package errors

type Unwrapper interface {
	error
	Unwrap() error
}

type Flusher interface {
	Flush() error
}
