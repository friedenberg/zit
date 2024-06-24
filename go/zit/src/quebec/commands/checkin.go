package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

type Checkin struct {
	Kasten     kennung.Kasten
	Delete     bool
	IgnoreAkte bool
}

func init() {
	registerCommandWithQuery(
		"checkin",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Checkin{}

			f.Var(&c.Kasten, "kasten", "none or Chrome")

			f.BoolVar(&c.Delete, "delete", false, "the checked-out file")

			f.BoolVar(
				&c.IgnoreAkte,
				"ignore-akte",
				false,
				"do not change the akte",
			)

			return c
		},
	)
}

func (c Checkin) DefaultGattungen() kennung.Gattung {
	return kennung.MakeGattung()
}

func (c Checkin) RunWithQuery(
	u *umwelt.Umwelt,
	ms *query.Group,
) (err error) {
	op := user_ops.Checkin{
		Delete: c.Delete,
	}

	if err = op.Run(
		u,
		query.GroupWithKasten{
			Group:  ms,
			Kasten: c.Kasten,
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
