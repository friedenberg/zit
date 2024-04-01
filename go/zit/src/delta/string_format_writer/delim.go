package string_format_writer

import (
	"bufio"
	"io"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
)

func MakeDelim[T any](
	delim string,
	w1 schnittstellen.WriterAndStringWriter,
	f schnittstellen.StringFormatWriter[T],
) func(T) error {
	w := bufio.NewWriter(w1)

	return iter.MakeSyncSerializer(
		func(e T) (err error) {
			errors.TodoP3("modify flushing behavior based on w1 being a TTY")
			defer errors.DeferredFlusher(&err, w)

			if _, err = f.WriteStringFormat(w, e); err != nil {
				err = errors.Wrap(err)
				return
			}

			if _, err = io.WriteString(w, delim); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	)
}
