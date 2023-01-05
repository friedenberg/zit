package format

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/collections"
)

func Write(
	w io.Writer,
	wffs ...WriterFunc,
) (n int64, err error) {
	for _, wf := range wffs {
		var n1 int64

		if n1, err = wf(w); err != nil {
			err = errors.Wrap(err)
			return
		}

		n += n1
	}

	return
}

// TODO rename
func MakeWriterTo[T any](
	w io.Writer,
	wf FormatWriterFunc[T],
) func(*T) error {
	return func(e *T) (err error) {
		if _, err = wf(w, e); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func MakeWriterToWithNewLines[T any](
	w io.Writer,
	wf FormatWriterFunc[T],
) func(*T) error {
	// w := bufio.NewWriter(w1)

	return collections.MakeSyncSerializer(
		func(e *T) (err error) {
			//TODO modify flushing behavior based on w1 being a TTY
			// defer errors.Deferred(&err, w.Flush)

			if _, err = wf(w, e); err != nil {
				err = errors.Wrap(err)
				return
			}

			if _, err = io.WriteString(w, "\n"); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	)
}
