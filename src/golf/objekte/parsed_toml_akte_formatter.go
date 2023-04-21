package objekte

import (
	"bufio"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/toml"
)

type ParsedAkteTomlFormatter[O Objekte[O]] struct{}

func (_ ParsedAkteTomlFormatter[O]) FormatParsedAkte(
	w1 io.Writer,
	t O,
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
