package etikett

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/toml"
)

// TODO-P2 merge into virtual objekte handling
func WriteObjekteToText(w io.Writer, t *Objekte) (n int64, err error) {
	enc := toml.NewEncoder(w)

	if err = enc.Encode(&t.Akte); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
