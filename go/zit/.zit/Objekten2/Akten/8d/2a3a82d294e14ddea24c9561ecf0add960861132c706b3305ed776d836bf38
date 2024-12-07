package repo_blobs

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/toml"
)

type FormatterBlobTextToml struct{}

func (FormatterBlobTextToml) Format(
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
