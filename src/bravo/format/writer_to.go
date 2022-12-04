package format

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type writerTo[T any] struct {
	wf FuncWriterElement[T]
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
	wf FuncWriterElement[T],
	e *T,
) *writerTo[T] {
	return &writerTo[T]{
		wf: wf,
		e:  e,
	}
}
