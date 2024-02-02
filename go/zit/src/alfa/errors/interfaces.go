package errors

import "github.com/friedenberg/zit/src/alfa/schnittstellen"

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

type FlusherWithLogger interface {
	Flush(schnittstellen.FuncIter[string]) error
}
