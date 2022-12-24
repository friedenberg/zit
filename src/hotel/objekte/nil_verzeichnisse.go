package objekte

import "github.com/friedenberg/zit/src/charlie/gattung"

type NilVerzeichnisse[T gattung.Objekte[T]] struct{}

func (_ NilVerzeichnisse[T]) ResetWithObjekte(o *T) {}
