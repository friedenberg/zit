package age

import (
	"io"
	"os"

	"filippo.io/age"
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
)

func Generate(basePath string) (a *Age, err error) {
	var i *X25519Identity

	if i, err = age.GenerateX25519Identity(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var f *os.File

	if f, err = files.CreateExclusiveWriteOnly(basePath); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, f.Close)

	if _, err = io.WriteString(f, i.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	a = &Age{
		recipients: []Recipient{i.Recipient()},
		identities: []Identity{i},
	}

	return
}
