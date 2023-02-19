package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/juliett/cwd"
	"github.com/friedenberg/zit/src/lima/zettel_checked_out"
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
	fds := kennung.MakeMutableFDSet()

	if pz, err = cwd.MakeCwdFilesMetaSet(
		u.Konfig(),
		u.Standort().Cwd(),
		ms,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	// ptl := u.PrinterTransactedLike()

	u.Lock()
	defer errors.Deferred(&err, u.Unlock)

	if err = u.StoreWorkingDirectory().ReadFiles(
		pz,
		func(co objekte.CheckedOutLike) (err error) {
			// var tl objekte.TransactedLike

			switch aco := co.(type) {
			case zettel_checked_out.Zettel:
				if _, err = u.StoreObjekten().Zettel().Update2(
					aco,
				); err != nil {
					err = errors.Wrap(err)
					return
				}

			case *zettel_checked_out.Zettel:
				if _, err = u.StoreObjekten().Zettel().Update2(
					*aco,
				); err != nil {
					err = errors.Wrap(err)
					return
				}

			default:
				errors.Todo("implement")
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

	if !c.Delete {
		return
	}

	deleteOp := user_ops.DeleteCheckout{
		Umwelt: u,
	}

	if err = deleteOp.Run(fds); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
