package pool

import (
	"bufio"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

var (
	bufioReader interfaces.Pool[bufio.Reader, *bufio.Reader]
	bufioWriter interfaces.Pool[bufio.Writer, *bufio.Writer]
)

func init() {
	bufioReader = MakePool[bufio.Reader, *bufio.Reader](nil, nil)
	bufioWriter = MakePool[bufio.Writer](nil, nil)
}

func GetBufioReader() interfaces.Pool[bufio.Reader, *bufio.Reader] {
	return bufioReader
}

func GetBufioWriter() interfaces.Pool[bufio.Writer, *bufio.Writer] {
	return bufioWriter
}
