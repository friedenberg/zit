package ohio

import (
	"bufio"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
)

var poolBufioWriter interfaces.Pool[bufio.Writer, *bufio.Writer]

func init() {
	poolBufioWriter = pool.MakePool[bufio.Writer, *bufio.Writer](nil, nil)
}

func GetPoolBufioWriter() interfaces.Pool[bufio.Writer, *bufio.Writer] {
	return poolBufioWriter
}
