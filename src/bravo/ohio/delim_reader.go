package ohio

import (
	"bufio"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
)

var delimReaderPool schnittstellen.Pool[delimReader, *delimReader]

func init() {
	delimReaderPool = collections.MakePool[delimReader, *delimReader]()
}

func PutDelimReader(dr *delimReader) {
	delimReaderPool.Put(dr)
}

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
) (dr *delimReader) {
	dr = delimReaderPool.Get()
	dr.br.Reset(r)
	dr.delim = delim

	return
}

func (lr *delimReader) N() int64 {
	return lr.n
}

func (lr *delimReader) Segments() int64 {
	return lr.segments
}

func (lr *delimReader) IsEOF() bool {
	return lr.eof
}

func (lr *delimReader) ResetWith(dr delimReader) {
	lr.Reset()
	lr.delim = dr.delim
}

func (lr *delimReader) Reset() {
	if lr.br == nil {
		lr.br = bufio.NewReader(nil)
	} else {
		lr.br.Reset(nil)
	}

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
