package delim_io

import (
	"fmt"
	"io"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
)

// Copies each `delim` suffixed segment from src to dst, and for each segment,
// adds the passed in prefix string.
//
// Useful for taking a Reader and adding a prefix for every line, like how `git`
// shows `remote: <line>` for all remote stderr output.
// TODO extract into an io.Writer-like object
func CopyWithPrefixOnDelim(
	delim byte,
	prefix string,
	dst ui.Printer,
	src io.Reader,
	includeLineNo bool,
) (n int64, err error) {
	br := pool.GetBufioReader().Get()
	defer pool.GetBufioReader().Put(br)
	br.Reset(src)

	var (
		eof    bool
		lineNo int
	)

	var sb strings.Builder

	for !eof {
		var rawLine string

		rawLine, err = br.ReadString(delim)
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

		sb.WriteString(prefix)
		fmt.Fprint(&sb, ":")

		if includeLineNo {
			fmt.Fprintf(&sb, "%d:", lineNo)
		}

		fmt.Fprint(&sb, " ")
		// fmt.Fprint(bw, "\t")

		sb.WriteString(strings.TrimSuffix(rawLine, string([]byte{delim})))

		dst.Print(sb.String())
		sb.Reset()

		lineNo++
	}

	return
}
