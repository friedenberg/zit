package gattung

import (
	"flag"
	"fmt"

	"github.com/friedenberg/zit/src/bravo/sha_core"
)

type Equatable[T any] interface {
	Equals(*T) bool
}

type Resetable[T any] interface {
	Reset(*T)
}

type Element interface {
}

type ElementPtr[T Element] interface {
	*T
}

type Keyer[T Element, T1 ElementPtr[T]] interface {
	Key(T1) string
}

type ValueElement interface {
	fmt.Stringer
}

type ValueElementPtr[T any] interface {
	*T
	flag.Value
}

type Identifier[T any] interface {
	Gattung() Gattung
	ValueElement
	Equatable[T]
}

type IdentifierPtr[T any] interface {
	ValueElementPtr[T]
	Reset(*T)
}

type Objekte interface {
	Gattung() Gattung
	AkteSha() sha_core.Sha
}

type ObjektePtr[T any] interface {
	*T
	Equatable[T]
	Resetable[T]

	SetAkteSha(sha_core.Sha)
}

type Stored interface {
	Gattung() Gattung

	AkteSha() sha_core.Sha
	SetAkteSha(sha_core.Sha)

	SetObjekteSha(AkteReaderFactory, string) error
	ObjekteSha() sha_core.Sha
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
