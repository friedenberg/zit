package standort

import (
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/debug"
)

type Options struct {
	BasePath string
	Debug    debug.Options
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
