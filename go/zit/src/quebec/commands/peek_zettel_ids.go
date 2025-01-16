package commands

import (
	"sort"
	"strconv"

	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	registerCommand("peek-zettel-ids", &PeekZettelIds{})
}

type PeekZettelIds struct {
	command_components.LocalWorkingCopy
}

func (cmd PeekZettelIds) Run(dep command.Dep) {
	args := dep.Args()

	n := 0

	if len(args) > 0 {
		{
			var err error

			if n, err = strconv.Atoi(args[0]); err != nil {
				dep.CancelWithErrorf("expected int but got %s", args[0])
			}
		}
	}

	localWorkingCopy := cmd.MakeLocalWorkingCopy(dep)

	var hs []*ids.ZettelId

	{
		var err error
		if hs, err = localWorkingCopy.GetStore().GetZettelIdIndex().PeekZettelIds(
			n,
		); err != nil {
			localWorkingCopy.CancelWithError(err)
		}
	}

	sort.Slice(
		hs,
		func(i, j int) bool {
			return hs[i].String() < hs[j].String()
		},
	)

	for i, h := range hs {
		localWorkingCopy.GetUI().Printf("%d: %s", i, h)
	}

	return
}
