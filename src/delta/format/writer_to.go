package format

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type writerTo[T any] struct {
	wf schnittstellen.FuncWriterElement[T]
	e  *T
}

func (wt *writerTo[T]) WriteTo(w io.Writer) (n int64, err error) {
	if n, err = wt.wf(w, wt.e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MakeWriterTo2[T any](
	wf schnittstellen.FuncWriterElement[T],
	e *T,
) *writerTo[T] {
	return &writerTo[T]{
		wf: wf,
		e:  e,
	}
}

type writerToInterface[T any] struct {
	wf schnittstellen.FuncWriterElementInterface[T]
	e  T
}

func (wt writerToInterface[T]) WriteTo(w io.Writer) (n int64, err error) {
	if n, err = wt.wf(w, wt.e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MakeWriterToInterface[T any](
	wf schnittstellen.FuncWriterElementInterface[T],
	e T,
) writerToInterface[T] {
	return writerToInterface[T]{
		wf: wf,
		e:  e,
	}
}
