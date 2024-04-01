package kasten_akte

import (
	"io"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/toml"
)

type FormatterAkteTextToml struct{}

func (_ FormatterAkteTextToml) Format(
	w io.Writer,
	t *V0,
) (n int64, err error) {
	enc := toml.NewEncoder(w)

	if err = enc.Encode(&t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
