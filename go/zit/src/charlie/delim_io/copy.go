package delim_io

import (
	"bytes"
	"fmt"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
)

// Copies each `delim` suffixed segment from src to dst, and for each segment,
// adds the passed in prefix string.
//
// Useful for taking a Reader and adding a prefix for every line, like how `git`
// shows `remote: <line>` for all remote stderr output.
func CopyWithPrefixOnDelim(
	delim byte,
	prefix string,
	dst io.Writer,
	src io.Reader,
	includeLineNo bool,
) (n int64, err error) {
	br := pool.GetBufioReader().Get()
	defer pool.GetBufioReader().Put(br)
	br.Reset(src)

	bw := pool.GetBufioWriter().Get()
	defer pool.GetBufioWriter().Put(bw)
	defer errors.DeferredFlusher(&err, bw)
	bw.Reset(dst)

	var (
		eof    bool
		lineNo int
	)

	for !eof {
		var rawLine []byte

		rawLine, err = br.ReadBytes(delim)
		n1 := len(rawLine)
		n += int64(n1)

		if err != nil && !errors.IsEOF(err) {
			err = errors.Wrap(err)
			return
		}

		if errors.IsEOF(err) {
			eof = true
			err = nil

			if n1 == 0 {
				break
			}
		}

		bw.WriteString(prefix)
		fmt.Fprint(bw, ":")

		if includeLineNo {
			fmt.Fprintf(bw, "%d:", lineNo)
		}

		fmt.Fprint(bw, " ")
		// fmt.Fprint(bw, "\t")

		bw.Write(bytes.TrimSuffix(rawLine, []byte{delim}))
		bw.WriteByte(delim)

		lineNo++
	}

	return
}
