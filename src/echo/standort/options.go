package standort

import (
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type Options struct {
	BasePath string
	cwd      string
}

func (o *Options) Validate() (err error) {
	if o.cwd, err = os.Getwd(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if o.BasePath == "" {
		o.BasePath = o.cwd
	}

	return
}
