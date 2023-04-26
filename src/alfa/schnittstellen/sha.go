package schnittstellen

import "io"

type Sha interface {
	// TODO-P3
	// GetHashBytes() []byte
	ValueLike
	Korper
	GetShaString() string
	EqualsSha(Sha) bool
	IsNull() bool
	ShaGetter
}

type ShaGetter interface {
	GetSha() Sha
}

type ShaReadCloser interface {
	io.WriterTo
	io.ReadCloser
	Sha() Sha
}

type ShaWriteCloser interface {
	io.ReaderFrom
	io.WriteCloser
	Sha() Sha
}

type (
	FuncShaReadCloser  func(Sha) (ShaReadCloser, error)
	FuncShaWriteCloser func(Sha) (ShaWriteCloser, error)
)
