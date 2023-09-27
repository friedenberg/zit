package ohio

import (
	"bufio"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/pool"
)

var poolBufioWriter schnittstellen.Pool[bufio.Writer, *bufio.Writer]

func init() {
	poolBufioWriter = pool.MakePool[bufio.Writer, *bufio.Writer](nil, nil)
}

func GetPoolBufioWriter() schnittstellen.Pool[bufio.Writer, *bufio.Writer] {
	return poolBufioWriter
}
