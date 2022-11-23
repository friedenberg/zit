package objekte_stored

import "github.com/friedenberg/zit/src/charlie/sha"

type Identifier interface {
}

type Named[T any, T1 ObjektePtr[T], T2 Identifier] struct {
	Stored  Stored[T, T1]
	Sha     sha.Sha
	Kennung Identifier
}
