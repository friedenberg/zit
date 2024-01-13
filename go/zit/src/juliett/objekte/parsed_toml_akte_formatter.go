package objekte

import (
	"bufio"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/alfa/toml"
)

type ParsedAkteTomlFormatter[O schnittstellen.Akte[O], OPtr schnittstellen.AktePtr[O]] struct{}

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
