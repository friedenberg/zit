package commands

import (
	"flag"
	"sort"
	"strconv"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type PeekZettelIds struct{}

func init() {
	registerCommand(
		"peek-zettel-ids",
		func(f *flag.FlagSet) Command {
			c := &PeekZettelIds{}

			return c
		},
	)
}

func (c PeekZettelIds) Run(store *env.Local, args ...string) (err error) {
	n := 0

	if len(args) > 0 {
		if n, err = strconv.Atoi(args[0]); err != nil {
			err = errors.Errorf("expected int but got %s", args[0])
			return
		}
	}

	var hs []*ids.ZettelId

	if hs, err = store.GetStore().GetZettelIdIndex().PeekZettelIds(
		n,
	); err != nil {
		err = errors.Wrap(err)
		return
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
