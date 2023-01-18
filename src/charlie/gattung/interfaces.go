package gattung

import (
	"fmt"

	"github.com/friedenberg/zit/src/schnittstellen"
)

type FuncAbbrId func(IdLike) (string, error)
type FuncAbbrIdMitKorper func(IdMitKorper) (string, error)

type ShaLike = schnittstellen.Sha

type IdLike interface {
	fmt.Stringer
}

type IdMitKorper interface {
	IdLike
	Kopf() string
	Schwanz() string
}

type Resetable[T any] interface {
	Reset(*T)
}

type Resetter[T any] interface {
	schnittstellen.Ptr[T]
	Reset2()
}

type ResetWither[T any, TPtr schnittstellen.Ptr[T]] interface {
	ResetWith(TPtr)
}

type Keyer[T any, T1 schnittstellen.Ptr[T]] interface {
	Key(T1) string
}

//   ___    _            _   _  __ _
//  |_ _|__| | ___ _ __ | |_(_)/ _(_) ___ _ __
//   | |/ _` |/ _ \ '_ \| __| | |_| |/ _ \ '__|
//   | | (_| |  __/ | | | |_| |  _| |  __/ |
//  |___\__,_|\___|_| |_|\__|_|_| |_|\___|_|
//

type IdentifierLike interface {
	schnittstellen.GattungGetter
	IdLike
}

type Id[T schnittstellen.Value] interface {
	schnittstellen.Equatable[T]
	fmt.Stringer
}

type IdPtr[T schnittstellen.Value] interface {
	Id[T]
	schnittstellen.ValuePtr[T]
}

// TODO-P2 rename to ObjekteKennung
type Identifier[T any] interface {
	schnittstellen.GattungGetter
	schnittstellen.Equatable[T]
	IdentifierLike
}

type IdentifierPtr[T schnittstellen.Value] interface {
	schnittstellen.ValuePtr[T]
	Resetable[T]
}

//    ___  _     _      _    _
//   / _ \| |__ (_) ___| | _| |_ ___
//  | | | | '_ \| |/ _ \ |/ / __/ _ \
//  | |_| | |_) | |  __/   <| ||  __/
//   \___/|_.__// |\___|_|\_\\__\___|
//            |__/

type Objekte[T any] interface {
	schnittstellen.GattungGetter
	schnittstellen.Equatable[T]
	GetAkteSha() schnittstellen.Sha
}

type ObjektePtr[T any] interface {
	Objekte[T]
	schnittstellen.Ptr[T]
	Resetable[T]
	SetAkteSha(schnittstellen.Sha)
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

type VerzeichnissePtr[T any, T1 Objekte[T1]] interface {
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
	Resetable[T]
	StoredPtr
}
