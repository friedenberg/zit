package schnittstellen

import "io"

// TODO-P3 refactor into hash or checksum or content address and split korper
// out into context object
type ShaLike interface {
	// TODO-P3
	// GetHashBytes() []byte
	// ValueLike
	Korper
	GetShaString() string
	GetShaBytes() []byte
	EqualsSha(ShaLike) bool // TODO-P3 rename to EqualsShaLike
	IsNull() bool
	//TODO AssertEquals() error
	ShaGetter
}

type ShaGetter interface {
	GetShaLike() ShaLike
}

type ShaReadCloser interface {
  io.Seeker
	io.WriterTo
	io.ReadCloser
	GetShaLike() ShaLike
}

type ShaWriteCloser interface {
	io.ReaderFrom
	io.WriteCloser
  // io.WriterAt
	GetShaLike() ShaLike
}

type (
	FuncShaReadCloser  func(ShaLike) (ShaReadCloser, error)
	FuncShaWriteCloser func(ShaLike) (ShaWriteCloser, error)
)
