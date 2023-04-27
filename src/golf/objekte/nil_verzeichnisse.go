package objekte

import "github.com/friedenberg/zit/src/foxtrot/metadatei"

type NilVerzeichnisse[T Akte[T]] struct{}

func (_ NilVerzeichnisse[T]) ResetWithObjekteMetadateiGetter(
	_ T,
	_ metadatei.Getter,
) {
}

func (_ *NilVerzeichnisse[T]) ResetWith(_ NilVerzeichnisse[T]) {}
func (_ *NilVerzeichnisse[T]) Reset()                          {}
