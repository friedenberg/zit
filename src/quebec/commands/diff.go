package commands

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/kilo/cwd"
	"github.com/friedenberg/zit/src/kilo/zettel_external"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type Diff struct{}

func init() {
	registerCommandWithCwdQuery(
		"diff",
		func(f *flag.FlagSet) CommandWithCwdQuery {
			c := &Diff{}

			return c
		},
	)
}

func (c Diff) DefaultGattungen() gattungen.Set {
	return gattungen.MakeSet(gattung.TrueGattung()...)
}

func (c Diff) RunWithCwdQuery(
	u *umwelt.Umwelt,
	ms kennung.MetaSet,
	cwdFiles cwd.CwdFiles,
) (err error) {
	e := zettel_external.MakeFileEncoderJustOpen(
		u.StoreObjekten(),
		u.Konfig(),
	)

	if err = u.StoreWorkingDirectory().ReadFiles(
		cwdFiles,
		ms,
		iter.MakeChain(
			objekte.MakeFilterFromMetaSet(ms),
			func(co objekte.CheckedOutLike) (err error) {
				var zco *zettel.CheckedOut
				ok := false

				if zco, ok = co.(*zettel.CheckedOut); !ok {
					return
				}

				var pFifo string

				if pFifo, err = u.Standort().FifoPipe(); err != nil {
					err = errors.Wrap(err)
					return
				}

				wg := iter.MakeErrorWaitGroup()

				wg.DoAfter(
					func() (err error) {
						return os.Remove(pFifo)
					},
				)

				wg.Do(
					func() error {
						return e.EncodeObjekte(
							&zco.Internal.Objekte,
							pFifo,
							"",
						)
					},
				)

				cmd := exec.Command(
					"diff",
					"--color=always",
					"-u",
					"--label", fmt.Sprintf("%s@zettel", zco.Internal.Sku.Kennung),
					pFifo,
					co.GetExternal().GetObjekteFD().Path,
				)

				cmd.Stdout = u.Out()
				cmd.Stderr = u.Err()

				wg.Do(
					func() (err error) {
						if err = cmd.Run(); err != nil {
							if cmd.ProcessState.ExitCode() == 1 {
								err = nil
							} else {
								err = errors.Wrap(err)
							}

							return
						}

						return
					},
				)

				return wg.GetError()
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
