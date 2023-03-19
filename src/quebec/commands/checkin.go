package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/gattungen"
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
	registerCommandWithCwdQuery(
		"checkin",
		func(f *flag.FlagSet) CommandWithCwdQuery {
			c := &Checkin{}

			f.BoolVar(&c.Delete, "delete", false, "the checked-out file")
			f.BoolVar(&c.IgnoreAkte, "ignore-akte", false, "do not change the akte")

			return c
		},
	)
}

func (c Checkin) DefaultGattungen() gattungen.Set {
	return gattungen.MakeSet()
}

func (c Checkin) RunWithCwdQuery(
	u *umwelt.Umwelt,
	ms kennung.MetaSet,
	pz cwd.CwdFiles,
) (err error) {
	op := user_ops.Checkin{
		Delete: c.Delete,
	}

	if err = op.Run(u, ms, pz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
