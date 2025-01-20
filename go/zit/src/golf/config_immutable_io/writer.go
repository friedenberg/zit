package config_immutable_io

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/triple_hyphen_io"
)

type Writer struct {
	*ConfigLoaded
}

func (writer Writer) WriteTo(w io.Writer) (n int64, err error) {
	thw := triple_hyphen_io.Writer{
		Metadata: metadata{ConfigLoaded: writer.ConfigLoaded},
		Blob:     writer.ConfigLoaded,
	}

	if n, err = thw.WriteTo(w); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
