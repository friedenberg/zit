package gattung

import (
	"github.com/friedenberg/zit/src/schnittstellen"
)

type Keyer[T any, T1 schnittstellen.Ptr[T]] interface {
	Key(T1) string
}

//   _____                               _           _
//  |_   _| __ __ _ _ __  ___  __ _  ___| |_ ___  __| |
//    | || '__/ _` | '_ \/ __|/ _` |/ __| __/ _ \/ _` |
//    | || | | (_| | | | \__ \ (_| | (__| ||  __/ (_| |
//    |_||_|  \__,_|_| |_|___/\__,_|\___|\__\___|\__,_|
//

type Transacted[T any] interface {
	schnittstellen.Stored
	GetKennungString() string
}

type TransactedPtr[T any] interface {
	Transacted[T]
	schnittstellen.Ptr[T]
	schnittstellen.Resetable[T]
	schnittstellen.StoredPtr
}
