package commands

import (
	"bufio"
	"flag"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/inventory_list_fmt"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
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
	u *env.Env,
	qg *query.Group,
) (err error) {
	conflicted := sku.MakeCheckedOutLikeMutableSet()

	if err = u.GetStore().QueryCheckedOut(
		qg,
		func(co sku.CheckedOutLike) (err error) {
			if co.GetState() != checked_out_state.Conflicted {
				return
			}

			if err = conflicted.Add(co.CloneCheckedOutLike()); err != nil {
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

	for col := range conflicted.All() {
		cofs := col.(*sku.CheckedOut)

		tm := sku.Conflicted{
			CheckedOutLike: col.CloneCheckedOutLike(),
		}

		var conflict *fd.FD

		if conflict, err = u.GetStore().GetCwdFiles().GetConflictOrError(&cofs.External); err != nil {
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

		s := inventory_list_fmt.MakeScanner(
			br,
			object_inventory_format.FormatForVersion(u.GetConfig().GetStoreVersion()),
			u.GetStore().GetObjekteFormatOptions(),
		)

		if err = tm.ReadConflictMarker(
			s,
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
