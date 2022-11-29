package gattung

import (
	"flag"
	"fmt"
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
	ValueElement
	Equatable[T]
}

type Identifier2[T any] interface {
	Gattung() Gattung
	ValueElement
	Equatable[T]
}

type IdentifierPtr[T any] interface {
	ValueElementPtr[T]
	Reset(*T)
}
