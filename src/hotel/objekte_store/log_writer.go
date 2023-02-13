package objekte_store

import "github.com/friedenberg/zit/src/alfa/schnittstellen"

type LogWriter[
	T any,
] struct {
	New, Updated, Unchanged, Archived schnittstellen.FuncIter[T]
}

func (l LogWriter[T]) NewOrUpdated(err error) schnittstellen.FuncIter[T] {
	if IsNotFound(err) {
		return l.New
	} else {
		return l.Updated
	}
}
