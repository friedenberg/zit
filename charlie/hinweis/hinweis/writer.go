package hinweis

import "io"

type writer struct {
	basePath string
}

func MakeWriter(basePath string) writer {
	return writer{
		basePath: basePath,
	}
}

func (w writer) WriteObjekte(s _Sha, out io.Writer) (err error) {
	return
}
