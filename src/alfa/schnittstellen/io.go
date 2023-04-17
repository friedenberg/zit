package schnittstellen

import "io"

type FuncReader func(io.Reader) (int64, error)

type (
	FuncReaderFormat[T any]  func(io.Reader, *T) (int64, error)
	FuncWriterElement[T any] func(io.Writer, *T) (int64, error)

	FuncReaderFormatInterface[T any]  func(io.Reader, T) (int64, error)
	FuncWriterElementInterface[T any] func(io.Writer, T) (int64, error)
)

type (
	FuncWriter              func(io.Writer) (int64, error)
	FuncWriterFormat[T any] func(io.Writer, T) (int64, error)
)
