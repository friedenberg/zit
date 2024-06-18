package standort

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/debug"
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
