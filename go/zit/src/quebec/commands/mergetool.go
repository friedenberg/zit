package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/juliett/to_merge"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

type Mergetool struct{}

func init() {
	registerCommandWithQuery(
		"merge-tool",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Mergetool{}

			return c
		},
	)
}

func (c Mergetool) DefaultGattungen() kennung.Gattung {
	return kennung.MakeGattung(gattung.TrueGattung()...)
}

func (c Mergetool) RunWithQuery(
	u *umwelt.Umwelt,
	qg *query.Group,
) (err error) {
	p := []string{}

	if err = u.GetStore().ReadFiles(
		qg,
		iter.MakeChain(
			func(co *sku.CheckedOut) (err error) {
				if co.State != checked_out_state.StateConflicted {
					return iter.MakeErrStopIteration()
				}

				return
			},
			func(co *sku.CheckedOut) (err error) {
				p = append(p, co.External.FDs.MakeConflictMarker())
				return
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	u.Lock()
	defer errors.Deferred(&err, u.Unlock)

	if len(p) == 0 {
		// TODO-P2 return status 1
		ui.Err().Printf("nothing to merge")
		return
	}

	for _, p1 := range p {
		tm := to_merge.Sku{
			ConflictMarkerPath: p1,
		}

		if err = tm.ReadConflictMarker(
			u.GetKonfig().GetStoreVersion(),
			u.GetStore().GetObjekteFormatOptions(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = u.GetStore().RunMergeTool(
			tm,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
