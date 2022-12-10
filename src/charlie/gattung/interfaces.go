package gattung

import (
	"flag"
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/bravo/sha_core"
)

type Equatable[T any] interface {
	Equals(*T) bool
}

type Resetable[T any] interface {
	Reset(*T)
}

type Element interface{}

type ElementPtr[T Element] interface {
	*T
}

type Keyer[T Element, T1 ElementPtr[T]] interface {
	Key(T1) string
}

type ValueElement interface {
	Element
	fmt.Stringer
}

type ValueElementPtr[T ValueElement] interface {
	ElementPtr[T]
	flag.Value
}

type IdentifierLike interface {
	Gattung() Gattung
	fmt.Stringer
}

type Identifier[T any] interface {
	IdentifierLike
	Gattung() Gattung
	ValueElement
	Equatable[T]
}

type IdentifierPtr[T ValueElement] interface {
	ValueElementPtr[T]
	Resetable[T]
}

type Objekte[T any] interface {
	Gattung() Gattung
	AkteSha() sha_core.Sha
	Equatable[T]
}

type ObjektePtr[T Element] interface {
	ElementPtr[T]
	Resetable[T]
	SetAkteSha(sha_core.Sha)
}

// TODO-P2 split into Stored and StoredPtr
type Stored interface {
	Gattung() Gattung
	// Identifier() IdentifierLike

	AkteSha() sha_core.Sha
	SetAkteSha(sha_core.Sha)

	SetObjekteSha(AkteReaderFactory, string) error
	ObjekteSha() sha_core.Sha
}

type StoredPtr interface {
	Stored
	SetAkteSha(sha_core.Sha)
	SetObjekteSha(AkteReaderFactory, string) error
}

type AkteIOFactory interface {
	AkteReaderFactory
	AkteWriterFactory
}

type AkteReaderFactory interface {
	AkteReader(sha_core.Sha) (sha_core.ReadCloser, error)
}

type AkteWriterFactory interface {
	AkteWriter() (sha_core.WriteCloser, error)
}

type AkteIOFactoryFactory interface {
	AkteFactory(Gattung) AkteIOFactory
}

type FormatReader[T any] interface {
	ReadFormat(io.Reader, *T) (int64, error)
}

type FormatWriter[T any] interface {
	WriteFormat(io.Writer, *T) (int64, error)
}

type Formatter[T any] interface {
	FormatReader[T]
	FormatWriter[T]
}
