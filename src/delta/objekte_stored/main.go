package objekte_stored

import (
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/bravo/gattung"
)

type Equatable[T any] interface {
	Equals(*T) bool
}

type Objekte interface {
	Gattung() gattung.Gattung
}

type ObjektePtr[T any] interface {
	*T
	Equatable[T]
	Reset(*T)
}

type Identifier[T any] interface {
	collections.ValueElement
	Equatable[T]
}

type IdentifierPtr[T any] interface {
	collections.ValueElementPtr[T]
	Reset(*T)
}
