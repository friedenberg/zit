package typ

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/toml"
)

type FormatterAkteTextToml struct{}

func (_ FormatterAkteTextToml) Format(
	w io.Writer,
	t *Objekte,
) (n int64, err error) {
	enc := toml.NewEncoder(w)

	if err = enc.Encode(&t.Akte); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
