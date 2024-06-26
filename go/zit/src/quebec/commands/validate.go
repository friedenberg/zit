package commands

import (
	"flag"
	"sync/atomic"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

type Validate struct{}

func init() {
	registerCommandWithQuery(
		"validate",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Validate{}

			return c
		},
	)
}

func (c Validate) ModifyBuilder(b *query.Builder) {
	b.WithDefaultGattungen(kennung.MakeGattung(gattung.Zettel)).
		WithDoNotMatchEmpty()
}

func (c Validate) RunWithQuery(
	u *umwelt.Umwelt,
	qg *query.Group,
) (err error) {
	var failureCount atomic.Int32

	if err = qg.GetExplicitCwdFDs().Each(
		func(f *fd.FD) (err error) {
			var h kennung.Hinweis

			if h, err = kennung.GetHinweis(f, true); err != nil {
				err = errors.Wrap(err)
				return
			}

			t := &store_fs.KennungFDPair{}

			if err = t.Kennung.SetWithKennung(h); err != nil {
				err = errors.Wrap(err)
				return
			}

			t.FDs.Objekte.ResetWith(f)

			if _, err = u.GetStore().GetCwdFiles().ReadCheckedOutFromKennungFDPair(
				store.ObjekteOptions{
					Mode: objekte_mode.ModeUpdateTai,
				},
				t,
			); err != nil {
				failureCount.Add(1)
				err = errors.Wrapf(err, "File: %q", f)
				ui.Err().Print(err)
				err = nil
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if failureCount.Load() > 0 {
		err = errors.Normalf("")
		return
	}

	return
}
