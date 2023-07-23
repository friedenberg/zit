package ohio

import (
	"bufio"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type delimReader struct {
	delim    byte
	br       *bufio.Reader
	n        int64
	segments int64
	eof      bool
}

func MakeDelimReader(
	delim byte,
	r io.Reader,
) delimReader {
	return delimReader{
		delim: delim,
		br:    bufio.NewReader(r),
	}
}

func (lr delimReader) N() int64 {
	return lr.n
}

func (lr delimReader) Segments() int64 {
	return lr.segments
}

func (lr delimReader) IsEOF() bool {
	return lr.eof
}

func (lr *delimReader) Reset(r io.Reader) {
	lr.br.Reset(r)
	lr.n = 0
	lr.segments = 0
	lr.eof = false
}

// Not safe for parallel use
func (lr *delimReader) ReadOneString() (str string, err error) {
	if lr.eof {
		err = io.EOF
		return
	}

	var rawLine string

	rawLine, err = lr.br.ReadString(lr.delim)
	n1 := len(rawLine)
	lr.n += int64(n1)

	if err != nil && !errors.IsEOF(err) {
		err = errors.Wrap(err)
		return
	}

	if errors.IsEOF(err) {
		lr.eof = true
	}

	str = strings.TrimSuffix(rawLine, string([]byte{lr.delim}))

	lr.segments++

	return
}

// Not safe for parallel use
func (lr *delimReader) ReadOneKeyValue(
	sep string,
) (key, val string, err error) {
	str, err := lr.ReadOneString()
	if err != nil {
		err = errors.Wrap(err)
		return
	}

	loc := strings.Index(str, sep)

	if loc == -1 {
		err = errors.Errorf(
			"expected at least one %q, but found none: %q",
			sep,
			str,
		)
		return
	}

	key = str[:loc]
	val = str[loc+1:]

	return
}
