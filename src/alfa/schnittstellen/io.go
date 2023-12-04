package schnittstellen

import (
	"io"
)

type FuncReader func(io.Reader) (int64, error)

type (
	FuncReaderFormat[T any]  func(io.Reader, *T) (int64, error)
	FuncWriterElement[T any] func(io.Writer, *T) (int64, error)

	// TODO-P3 switch to below
	FuncReaderFormatInterface[T any]  func(io.Reader, T) (int64, error)
	FuncReaderElementInterface[T any] func(io.Writer, T) (int64, error)
	FuncWriterElementInterface[T any] func(io.Writer, T) (int64, error)
)

type (
	WriterAndStringWriter interface {
		io.Writer
		io.StringWriter
	}

	FuncWriter              func(io.Writer) (int64, error)
	FuncWriterFormat[T any] func(io.Writer, T) (int64, error)

	BufferSlice interface {
		Len() int
		CutBufferSlice(delim byte) (before, after BufferSlice, ok bool)
		Compare([]byte) int
		Equal([]byte) bool
		io.WriterTo
		io.ReaderFrom
		Stringer
		Bytes() []byte
	}

	RingBuffer interface {
		GetRingBuffer() RingBuffer
		PeekReadableSlice() BufferSlice
		PeekReadableSliceUpto(b byte) (readable BufferSlice, ok bool, err error)
		AdvanceRead(n int)
	}

	StringFormatReader[T any] interface {
		ReadStringFormat(RingBuffer, T) (int64, error)
	}

	StringFormatWriter[T any] interface {
		WriteStringFormat(WriterAndStringWriter, T) (int64, error)
	}

	StringFormatReadWriter[T any] interface {
		StringFormatReader[T]
		StringFormatWriter[T]
	}

	FuncStringWriterFormat[T any] func(WriterAndStringWriter, T) (int64, error)
)
