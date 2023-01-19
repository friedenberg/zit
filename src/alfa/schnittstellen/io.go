package schnittstellen

import "io"

type FuncReader func(io.Reader) (int64, error)

type FuncReaderFormat[T any] func(io.Reader, *T) (int64, error)
type FuncWriterElement[T any] func(io.Writer, *T) (int64, error)

// TODO rename to Func-prefix
type FuncWriter func(io.Writer) (int64, error)
type FuncWriterFormat[T any] func(io.Writer, T) (int64, error)
