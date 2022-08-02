package hinweis

import (
	"io"

	"github.com/friedenberg/zit/bravo/sha"
)

type writer struct {
	basePath string
}

func MakeWriter(basePath string) writer {
	return writer{
		basePath: basePath,
	}
}

func (w writer) WriteObjekte(s sha.Sha, out io.Writer) (err error) {
	return
}
