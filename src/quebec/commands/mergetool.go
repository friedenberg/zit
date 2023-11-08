package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/checked_out_state"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/juliett/to_merge"
	"github.com/friedenberg/zit/src/oscar/umwelt"
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

func (c Mergetool) DefaultGattungen() gattungen.Set {
	return gattungen.MakeSet(gattung.TrueGattung()...)
}

func (c Mergetool) RunWithQuery(
	u *umwelt.Umwelt,
	ms matcher.Query,
) (err error) {
	p := []string{}

	if err = u.StoreObjekten().ReadFiles(
		matcher.MakeFuncReaderTransactedLikePtr(ms, u.StoreObjekten().Query),
		iter.MakeChain(
			matcher.MakeFilterFromQuery(ms),
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
		errors.Err().Printf("nothing to merge")
		return
	}

	for _, p1 := range p {
		tm := to_merge.Sku{
			ConflictMarkerPath: p1,
		}

		if err = tm.ReadConflictMarker(
			u.Konfig().GetStoreVersion(),
			u.StoreUtil().GetObjekteFormatOptions(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = u.StoreObjekten().RunMergeTool(
			tm,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
