package akten

import (
	"bufio"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/toml"
)

type ParsedAkteTomlFormatter[O interfaces.Blob[O], OPtr interfaces.BlobPtr[O]] struct{}

func (_ ParsedAkteTomlFormatter[O, OPtr]) FormatParsedAkte(
	w1 io.Writer,
	t OPtr,
) (n int64, err error) {
	w := bufio.NewWriter(w1)
	defer errors.DeferredFlusher(&err, w)

	enc := toml.NewEncoder(w)

	if err = enc.Encode(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
