package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/kilo/cwd"
	"github.com/friedenberg/zit/src/november/umwelt"
	"github.com/friedenberg/zit/src/oscar/user_ops"
)

type Clean struct{}

func init() {
	registerCommandWithCwdQuery(
		"clean",
		func(f *flag.FlagSet) CommandWithCwdQuery {
			c := &Clean{}

			return c
		},
	)
}

func (c Clean) RunWithCwdQuery(
	s *umwelt.Umwelt,
	possible cwd.CwdFiles,
) (err error) {
	fds := kennung.MakeMutableFDSet()

	for _, d := range possible.EmptyDirectories {
		fds.Add(d)
	}

	if err = s.StoreWorkingDirectory().ReadFiles(
		possible,
		func(co objekte.CheckedOutLike) (err error) {
			if co.GetState() != objekte.CheckedOutStateExistsAndSame {
				return
			}

			e := co.GetExternal()

			fds.Add(e.GetObjekteFD())
			fds.Add(e.GetAkteFD())

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	deleteOp := user_ops.DeleteCheckout{
		Umwelt: s,
	}

	if err = deleteOp.Run(fds); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
