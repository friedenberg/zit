package schnittstellen

import "io"

type Sha interface {
	//TODO
	// GetHashBytes() []byte
	Korper
	GetShaString() string
	Equals(Sha) bool
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

type FuncShaReadCloser func(Sha) (ShaReadCloser, error)
type FuncShaWriteCloser func(Sha) (ShaWriteCloser, error)
