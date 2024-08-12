package fs_home

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/debug"
)

type Options struct {
	BasePath             string
	Debug                debug.Options
	DryRun               bool
	PermitNoZitDirectory bool
	store_fs             string
}

func (o *Options) Validate() (err error) {
	if o.store_fs, err = os.Getwd(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
