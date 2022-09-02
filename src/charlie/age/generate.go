package age

import (
	"io/ioutil"
	"path"

	"filippo.io/age"
	"github.com/friedenberg/zit/src/alfa/errors"
)

func Generate(basePath string) (a *ages, err error) {
	var i *X25519Identity

	if i, err = age.GenerateX25519Identity(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = ioutil.WriteFile(path.Join(basePath), []byte(i.String()), 0755); err != nil {
		err = errors.Wrap(err)
		return
	}

	a = &ages{
		recipient: i.Recipient(),
		identity:  i,
	}

	return
}
