package gattung

import (
	"flag"
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/bravo/sha"
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

//   ___    _            _   _  __ _
//  |_ _|__| | ___ _ __ | |_(_)/ _(_) ___ _ __
//   | |/ _` |/ _ \ '_ \| __| | |_| |/ _ \ '__|
//   | | (_| |  __/ | | | |_| |  _| |  __/ |
//  |___\__,_|\___|_| |_|\__|_|_| |_|\___|_|
//

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

//    ___  _     _      _    _
//   / _ \| |__ (_) ___| | _| |_ ___
//  | | | | '_ \| |/ _ \ |/ / __/ _ \
//  | |_| | |_) | |  __/   <| ||  __/
//   \___/|_.__// |\___|_|\_\\__\___|
//            |__/

type Objekte[T any] interface {
	Gattung() Gattung
	Equatable[T]
	AkteSha() sha.Sha
}

type ObjektePtr[T Element] interface {
	Objekte[T]
	ElementPtr[T]
	Resetable[T]
	SetAkteSha(sha.Sha)
}

//   ____  _                     _
//  / ___|| |_ ___  _ __ ___  __| |
//  \___ \| __/ _ \| '__/ _ \/ _` |
//   ___) | || (_) | | |  __/ (_| |
//  |____/ \__\___/|_|  \___|\__,_|
//

type Stored interface {
	Gattung() Gattung
	//TODO-P4 add identifier
	// Identifier() IdentifierLike

	AkteSha() sha.Sha
	ObjekteSha() sha.Sha
}

type StoredPtr interface {
	Stored
	SetAkteSha(sha.Sha)
	SetObjekteSha(AkteReaderFactory, string) error
}

//  __     __                _      _           _
//  \ \   / /__ _ __ _______(_) ___| |__  _ __ (_)___ ___  ___
//   \ \ / / _ \ '__|_  / _ \ |/ __| '_ \| '_ \| / __/ __|/ _ \
//    \ V /  __/ |   / /  __/ | (__| | | | | | | \__ \__ \  __/
//     \_/ \___|_|  /___\___|_|\___|_| |_|_| |_|_|___/___/\___|
//

type Verzeichnisse[T any] interface {
}

type VerzeichnissePtr[T Element, T1 Objekte[T1]] interface {
	Verzeichnisse[T]
	ResetWithObjekte(*T1)
}

//   _____                               _           _
//  |_   _| __ __ _ _ __  ___  __ _  ___| |_ ___  __| |
//    | || '__/ _` | '_ \/ __|/ _` |/ __| __/ _ \/ _` |
//    | || | | (_| | | | \__ \ (_| | (__| ||  __/ (_| |
//    |_||_|  \__,_|_| |_|___/\__,_|\___|\__\___|\__,_|
//

type Transacted[T Element] interface {
	Equatable[T]
	Stored

	GetKennungString() string
}

type TransactedPtr[T Element] interface {
	Transacted[T]
	ElementPtr[T]
	Resetable[T]
	StoredPtr
}

//      _    _    _       ___ ___
//     / \  | | _| |_ ___|_ _/ _ \
//    / _ \ | |/ / __/ _ \| | | | |
//   / ___ \|   <| ||  __/| | |_| |
//  /_/   \_\_|\_\\__\___|___\___/
//

type AkteIOFactory interface {
	AkteReaderFactory
	AkteWriterFactory
}

type AkteReaderFactory interface {
	AkteReader(sha.Sha) (sha.ReadCloser, error)
}

type AkteWriterFactory interface {
	AkteWriter() (sha.WriteCloser, error)
}

type FuncReadCloser func(sha.Sha) (sha.ReadCloser, error)
type FuncWriteCloser func(sha.Sha) (sha.WriteCloser, error)

type AkteIOFactoryFactory interface {
	AkteFactory(Gattung) AkteIOFactory
}

//   _____                          _
//  |  ___|__  _ __ _ __ ___   __ _| |_
//  | |_ / _ \| '__| '_ ` _ \ / _` | __|
//  |  _| (_) | |  | | | | | | (_| | |_
//  |_|  \___/|_|  |_| |_| |_|\__,_|\__|
//

type FormatReader[T Element, T1 ElementPtr[T]] interface {
	ReadFormat(io.Reader, T1) (int64, error)
}

type FormatWriter[T Element, T1 ElementPtr[T]] interface {
	WriteFormat(io.Writer, T1) (int64, error)
}

type Parser[T Element, T1 ElementPtr[T]] interface {
	Parse(io.Reader, T1) (int64, error)
}

type Formatter[T Element, T1 ElementPtr[T]] interface {
	Format(io.Writer, T1) (int64, error)
}

type Format[T Element, T1 ElementPtr[T]] interface {
	Parser[T, T1]
	Formatter[T, T1]
}
