package schnittstellen

import "io"

type ShaLike interface {
	// TODO-P3
	// GetHashBytes() []byte
	ValueLike
	Korper
	GetShaString() string
	EqualsSha(ShaLike) bool
	IsNull() bool
	ShaGetter
}

type ShaGetter interface {
	GetSha() ShaLike
}

type ShaReadCloser interface {
	io.WriterTo
	io.ReadCloser
	GetShaLike() ShaLike
}

type ShaWriteCloser interface {
	io.ReaderFrom
	io.WriteCloser
	GetShaLike() ShaLike
}

type (
	FuncShaReadCloser  func(ShaLike) (ShaReadCloser, error)
	FuncShaWriteCloser func(ShaLike) (ShaWriteCloser, error)
)
