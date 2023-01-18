package objekte

import "github.com/friedenberg/zit/src/schnittstellen"

type NilVerzeichnisse[T schnittstellen.Objekte[T]] struct{}

func (_ NilVerzeichnisse[T]) ResetWithObjekte(o *T) {}
