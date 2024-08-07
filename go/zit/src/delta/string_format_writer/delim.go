package string_format_writer

import (
	"bufio"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
)

func MakeDelim[T any](
	delim string,
	w1 interfaces.WriterAndStringWriter,
	f interfaces.StringFormatWriter[T],
) func(T) error {
	w := bufio.NewWriter(w1)

	return iter.MakeSyncSerializer(
		func(e T) (err error) {
			ui.TodoP3("modify flushing behavior based on w1 being a TTY")
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
