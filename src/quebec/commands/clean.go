package commands

import (
	"flag"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections_ptr"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/checked_out_state"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/lima/cwd"
	"github.com/friedenberg/zit/src/oscar/umwelt"
	"github.com/friedenberg/zit/src/oscar/user_ops"
)

type Clean struct {
	force bool
}

func init() {
	registerCommandWithCwdQuery(
		"clean",
		func(f *flag.FlagSet) CommandWithCwdQuery {
			c := &Clean{}

			f.BoolVar(
				&c.force,
				"force",
				false,
				"remove objekten in working directory even if they have changes",
			)

			return c
		},
	)
}

func (c Clean) DefaultGattungen() gattungen.Set {
	return gattungen.MakeSet(gattung.TrueGattung()...)
}

func (c Clean) RunWithCwdQuery(
	s *umwelt.Umwelt,
	ms matcher.Query,
	possible *cwd.CwdFiles,
) (err error) {
	fds := collections_ptr.MakeMutableValueSet[kennung.FD, *kennung.FD](nil)
	l := &sync.Mutex{}

	for _, d := range possible.EmptyDirectories {
		fds.Add(d)
	}

	if err = s.StoreWorkingDirectory().ReadFiles(
		possible,
		ms,
		iter.MakeChain(
			objekte.MakeFilterFromMetaSet(ms),
			func(co *objekte.CheckedOut2) (err error) {
				if co.GetState() != checked_out_state.StateExistsAndSame && !c.force {
					return
				}

				e := co.GetExternalLike()

				l.Lock()
				defer l.Unlock()

				fds.Add(e.GetObjekteFD())
				fds.Add(e.GetAkteFD())

				return
			},
		),
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
