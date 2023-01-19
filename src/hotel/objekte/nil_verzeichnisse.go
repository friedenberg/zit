package objekte

import "github.com/friedenberg/zit/src/alfa/schnittstellen"

type NilVerzeichnisse[T schnittstellen.Objekte[T]] struct{}

func (_ NilVerzeichnisse[T]) ResetWithObjekte(o T) {}
