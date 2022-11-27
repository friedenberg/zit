package objekte

import (
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/metadatei_io"
)

type Objekte interface {
	Gattung() gattung.Gattung
}

type Objekte2 interface {
	Objekte
	AkteSha() sha.Sha
}

type Stored2 interface {
	Gattung() gattung.Gattung
	AkteSha() sha.Sha
	Sha() sha.Sha
	SetSha(metadatei_io.AkteReaderFactory, string) error
}

type ObjektePtr[T any] interface {
	*T
	collections.Equatable[T]
	collections.Resetable[T]
}

type Identifier[T any] interface {
	collections.ValueElement
	collections.Equatable[T]
}

type Identifier2[T any] interface {
	Gattung() gattung.Gattung
	collections.ValueElement
	collections.Equatable[T]
}

type IdentifierPtr[T any] interface {
	collections.ValueElementPtr[T]
	Reset(*T)
}
