package age

import (
	"io/ioutil"
	"path"
)

func Generate(basePath string) (a *age, err error) {
	var i *_AgeX25519Identity

	if i, err = _GenerateX25519Identity(); err != nil {
		err = _Error(err)
		return
	}

	if err = ioutil.WriteFile(path.Join(basePath), []byte(i.String()), 0755); err != nil {
		err = _Error(err)
		return
	}

	a = &age{
		recipient: i.Recipient(),
		identity:  i,
	}

	return
}
