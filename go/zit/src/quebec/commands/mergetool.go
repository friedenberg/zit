package commands

import (
	"bufio"
	"flag"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/objekte_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
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

func (c Mergetool) DefaultGattungen() ids.Genre {
	return ids.MakeGenre(gattung.TrueGattung()...)
}

func (c Mergetool) RunWithQuery(
	u *umwelt.Umwelt,
	qg *query.Group,
) (err error) {
	conflicted := sku.MakeCheckedOutLikeMutableSet()

	if err = u.GetStore().QueryCheckedOut(
		qg,
		func(co sku.CheckedOutLike) (err error) {
			if co.GetState() != checked_out_state.StateConflicted {
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
		// TODO-P2 return status 1
		ui.Err().Printf("nothing to merge")
		return
	}

	if err = conflicted.Each(
		func(col sku.CheckedOutLike) (err error) {
			cofs := col.(*store_fs.CheckedOut)
			p1 := cofs.External.FDs.MakeConflictMarker()

			tm := sku.Conflicted{
				CheckedOutLike: col.Clone(),
			}

			var f *os.File

			if f, err = files.Open(p1); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer errors.DeferredCloser(&err, f)

			br := bufio.NewReader(f)

			s := sku_fmt.MakeFormatBestandsaufnahmeScanner(
				br,
				objekte_format.FormatForVersion(u.GetKonfig().GetStoreVersion()),
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

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
