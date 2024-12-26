package commands

import (
	"flag"
	"sort"
	"strconv"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type PeekZettelIds struct{}

func init() {
	registerCommand(
		"peek-zettel-ids",
		func(f *flag.FlagSet) CommandWithRepo {
			c := &PeekZettelIds{}

			return c
		},
	)
}

func (c PeekZettelIds) RunWithRepo(repo *repo_local.Repo, args ...string) {
	n := 0

	if len(args) > 0 {
		{
			var err error

			if n, err = strconv.Atoi(args[0]); err != nil {
				repo.CancelWithError(errors.Errorf("expected int but got %s", args[0]))
				return
			}
		}
	}

	var hs []*ids.ZettelId

	{
		var err error
		if hs, err = repo.GetStore().GetZettelIdIndex().PeekZettelIds(
			n,
		); err != nil {
			repo.CancelWithError(err)
			return
		}
	}

	sort.Slice(
		hs,
		func(i, j int) bool {
			return hs[i].String() < hs[j].String()
		},
	)

	for i, h := range hs {
		ui.Out().Printf("%d: %s", i, h)
	}

	return
}
