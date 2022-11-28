package objekte

import (
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/objekte_format"
)

type Objekte = objekte_format.Objekte
type Objekte2 = objekte_format.Objekte2
type Stored2 = objekte_format.Stored2

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
