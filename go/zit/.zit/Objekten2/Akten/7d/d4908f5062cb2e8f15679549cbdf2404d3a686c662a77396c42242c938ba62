package commands

import (
	"bufio"
	"flag"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
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

func (c Mergetool) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.TrueGenre()...)
}

func (c Mergetool) RunWithQuery(
	u *env.Local,
	qg *query.Group,
) (err error) {
	conflicted := sku.MakeSkuTypeSetMutable()

	if err = u.GetStore().QuerySkuType(
		qg,
		func(co sku.SkuType) (err error) {
			if co.GetState() != checked_out_state.Conflicted {
				return
			}

			if err = conflicted.Add(co.Clone()); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	u.Lock()
	defer errors.Deferred(&err, u.Unlock)

	if conflicted.Len() == 0 {
		// TODO-P2 return status 1 and use Err
		ui.Out().Printf("nothing to merge")
		return
	}

	for co := range conflicted.All() {
		tm := sku.Conflicted{
			CheckedOut: co.Clone(),
		}

		var conflict *fd.FD

		if conflict, err = u.GetStore().GetStoreFS().GetConflictOrError(
			co.GetSkuExternal(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		var f *os.File

		if f, err = files.Open(conflict.GetPath()); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, f)

		br := bufio.NewReader(f)

		bs := u.GetStore().GetBlobStore().GetInventoryList()

		if err = tm.ReadConflictMarker(
			func(f interfaces.FuncIter[*sku.Transacted]) (err error) {
				if err = bs.StreamInventoryListBlobSkusFromReader(
					builtin_types.DefaultOrPanic(genres.InventoryList),
					br,
					f,
				); err != nil {
					err = errors.Wrap(err)
					return
				}

        return
			},
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
