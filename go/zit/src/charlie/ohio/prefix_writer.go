package ohio

import (
	"bufio"
	"fmt"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
)

var (
	poolPrefixOnDelimReader schnittstellen.Pool[bufio.Reader, *bufio.Reader]
	poolPrefixOnDelimWriter schnittstellen.Pool[bufio.Writer, *bufio.Writer]
)

func init() {
	poolPrefixOnDelimReader = pool.MakePool[bufio.Reader, *bufio.Reader](
		nil,
		nil,
	)
	poolPrefixOnDelimWriter = pool.MakePool[bufio.Writer, *bufio.Writer](
		nil,
		nil,
	)
}

func CopyWithPrefixOnDelim(
	delim byte,
	prefix string,
	dst io.Writer,
	src io.Reader,
) (n int64, err error) {
	br := poolPrefixOnDelimReader.Get()
	defer poolPrefixOnDelimReader.Put(br)
	br.Reset(src)

	bw := poolPrefixOnDelimWriter.Get()
	defer poolPrefixOnDelimWriter.Put(bw)
	defer errors.DeferredFlusher(&err, bw)
	bw.Reset(dst)

	var (
		eof    bool
		lineNo int
	)

	for {
		if eof {
			err = nil
			break
		}

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
		}

		bw.WriteString(prefix)
		bw.WriteString(fmt.Sprintf(":%d:\t", lineNo))
		bw.WriteString(rawLine)
		bw.WriteByte(delim)

		lineNo++
	}

	return
}
