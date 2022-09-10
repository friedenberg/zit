package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/india/store_working_directory"
	"github.com/friedenberg/zit/src/juliett/store_with_lock"
	"github.com/friedenberg/zit/src/juliett/umwelt"
)

type GetPossibleZettels struct {
	umwelt *umwelt.Umwelt
}

func NewGetPossibleZettels(u *umwelt.Umwelt) GetPossibleZettels {
	return GetPossibleZettels{
		umwelt: u,
	}
}

func (uo GetPossibleZettels) Run(
	store store_with_lock.Store,
) (result store_working_directory.CwdFiles, err error) {
	if result, err = store.StoreWorkingDirectory().GetPossibleZettels(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
