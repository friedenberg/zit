package ohio

import (
	"bufio"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

func WriteSeq[T any](
	w1 io.Writer,
	e T,
	seq ...schnittstellen.FuncWriterElementInterface[T],
) (n int64, err error) {
	w := bufio.NewWriter(w1)
	defer errors.DeferredFlusher(&err, w)

	var n1 int64

	for _, s := range seq {
		n1, err = s(w, e)
		n += n1

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

// TODO check performance of this
func WriteLine(w io.Writer, s string) (n int64, err error) {
	var n1 int

	if s != "" {
		n1, err = io.WriteString(w, s)

		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	n1, err = io.WriteString(w, "\n")

	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}