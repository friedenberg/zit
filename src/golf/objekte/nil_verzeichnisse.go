package objekte

type NilVerzeichnisse[T Objekte[T]] struct{}

func (_ NilVerzeichnisse[T]) ResetWithObjekte(o T)             {}
func (_ *NilVerzeichnisse[T]) ResetWith(o NilVerzeichnisse[T]) {}
func (_ *NilVerzeichnisse[T]) Reset()                          {}
