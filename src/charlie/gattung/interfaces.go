package gattung

import (
	"github.com/friedenberg/zit/src/schnittstellen"
)

type FuncAbbrId func(schnittstellen.Value) (string, error)
type FuncAbbrIdMitKorper func(schnittstellen.Korper) (string, error)

type ShaLike = schnittstellen.Sha

type Keyer[T any, T1 schnittstellen.Ptr[T]] interface {
	Key(T1) string
}

//   ____  _                     _
//  / ___|| |_ ___  _ __ ___  __| |
//  \___ \| __/ _ \| '__/ _ \/ _` |
//   ___) | || (_) | | |  __/ (_| |
//  |____/ \__\___/|_|  \___|\__,_|
//

type Stored interface {
	schnittstellen.GattungGetter
	GetAkteSha() schnittstellen.Sha
	GetObjekteSha() schnittstellen.Sha
}

type StoredPtr interface {
	Stored
	SetAkteSha(schnittstellen.Sha)
	SetObjekteSha(schnittstellen.AkteReaderFactory, string) error
}

//  __     __                _      _           _
//  \ \   / /__ _ __ _______(_) ___| |__  _ __ (_)___ ___  ___
//   \ \ / / _ \ '__|_  / _ \ |/ __| '_ \| '_ \| / __/ __|/ _ \
//    \ V /  __/ |   / /  __/ | (__| | | | | | | \__ \__ \  __/
//     \_/ \___|_|  /___\___|_|\___|_| |_|_| |_|_|___/___/\___|
//

type Verzeichnisse[T any] interface {
}

type VerzeichnissePtr[T any, T1 schnittstellen.Objekte[T1]] interface {
	schnittstellen.Ptr[T]
	Verzeichnisse[T]
	ResetWithObjekte(*T1)
}

//   _____                               _           _
//  |_   _| __ __ _ _ __  ___  __ _  ___| |_ ___  __| |
//    | || '__/ _` | '_ \/ __|/ _` |/ __| __/ _ \/ _` |
//    | || | | (_| | | | \__ \ (_| | (__| ||  __/ (_| |
//    |_||_|  \__,_|_| |_|___/\__,_|\___|\__\___|\__,_|
//

type Transacted[T any] interface {
	Stored
	GetKennungString() string
}

type TransactedPtr[T any] interface {
	Transacted[T]
	schnittstellen.Ptr[T]
	schnittstellen.Resetable[T]
	StoredPtr
}
