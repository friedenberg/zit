package format

import (
	"io"

	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
)

type readerFrom[T any] struct {
	rf schnittstellen.FuncReaderFormat[T]
	e  *T
}

func (rf readerFrom[T]) ReadFrom(r io.Reader) (n int64, err error) {
	return rf.rf(r, rf.e)
}

func MakeReaderFrom[T any](
	rf schnittstellen.FuncReaderFormat[T],
	e *T,
) io.ReaderFrom {
	return readerFrom[T]{
		rf: rf,
		e:  e,
	}
}

type readerFromInterface[T any] struct {
	rf schnittstellen.FuncReaderFormatInterface[T]
	e  T
}

func (rf readerFromInterface[T]) ReadFrom(r io.Reader) (n int64, err error) {
	return rf.rf(r, rf.e)
}

func MakeReaderFromInterface[T any](
	rf schnittstellen.FuncReaderFormatInterface[T],
	e T,
) io.ReaderFrom {
	return readerFromInterface[T]{
		rf: rf,
		e:  e,
	}
}
