package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/india/store_working_directory"
	"github.com/friedenberg/zit/src/juliett/umwelt"
)

type GetPossibleZettels struct {
	*umwelt.Umwelt
}

func NewGetPossibleZettels(u *umwelt.Umwelt) GetPossibleZettels {
	return GetPossibleZettels{
		Umwelt: u,
	}
}

func (uo GetPossibleZettels) Run() (result store_working_directory.CwdFiles, err error) {
	if result, err = uo.StoreWorkingDirectory().GetPossibleZettels(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
