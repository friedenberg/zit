package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/kilo/cwd"
	"github.com/friedenberg/zit/src/november/umwelt"
	"github.com/friedenberg/zit/src/oscar/user_ops"
)

type Checkin struct {
	Delete     bool
	IgnoreAkte bool
}

func init() {
	registerCommandWithQuery(
		"checkin",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Checkin{}

			f.BoolVar(&c.Delete, "delete", false, "the checked-out file")
			f.BoolVar(&c.IgnoreAkte, "ignore-akte", false, "do not change the akte")

			return c
		},
	)
}

func (c Checkin) RunWithQuery(
	u *umwelt.Umwelt,
	ms kennung.MetaSet,
) (err error) {
	var pz cwd.CwdFiles

	if pz, err = cwd.MakeCwdFilesMetaSet(
		u.Konfig(),
		u.Standort().Cwd(),
		ms,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	op := user_ops.Checkin{
		Delete: c.Delete,
	}

	if err = op.Run(u, pz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
