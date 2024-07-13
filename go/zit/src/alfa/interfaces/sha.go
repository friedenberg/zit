package interfaces

import "io"

// TODO-P3 refactor into hash or checksum or content address and split korper
// out into context object
type Sha interface {
	// TODO-P3
	// GetHashBytes() []byte
	// ValueLike
	StringerWithHeadAndTail
	GetShaString() string
	GetShaBytes() []byte
	EqualsSha(Sha) bool // TODO-P3 rename to EqualsShaLike
	IsNull() bool
	ShaGetter
}

type ShaGetter interface {
	GetShaLike() Sha
}

type ShaReadCloser interface {
	io.Seeker
	io.WriterTo
	io.ReadCloser
	GetShaLike() Sha
}

type ShaWriteCloser interface {
	io.ReaderFrom
	io.WriteCloser
	// io.WriterAt
	GetShaLike() Sha
}

type (
	FuncShaReadCloser  func(Sha) (ShaReadCloser, error)
	FuncShaWriteCloser func(Sha) (ShaWriteCloser, error)
)
