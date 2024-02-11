package ohio

import (
	"bufio"

	"code.linenisgreat.com/zit-go/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit-go/src/bravo/pool"
)

var poolBufioWriter schnittstellen.Pool[bufio.Writer, *bufio.Writer]

func init() {
	poolBufioWriter = pool.MakePool[bufio.Writer, *bufio.Writer](nil, nil)
}

func GetPoolBufioWriter() schnittstellen.Pool[bufio.Writer, *bufio.Writer] {
	return poolBufioWriter
}
