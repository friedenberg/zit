package commands

import (
	"flag"
	"fmt"
	"path"

	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type Deinit struct {
	Force bool
}

func init() {
	registerCommand(
		"deinit",
		func(f *flag.FlagSet) CommandWithRepo {
			c := &Deinit{}

			f.BoolVar(
				&c.Force,
				"force",
				false,
				"force deinit",
			)

			return c
		},
	)
}

func (c Deinit) RunWithRepo(u *repo_local.Repo, args ...string) {
	if !c.Force && !c.getPermission(u) {
		return
	}

	base := path.Join(u.GetRepoLayout().Dir())

	if err := files.SetAllowUserChangesRecursive(base); err != nil {
		u.CancelWithError(err)
	}

	if err := u.GetRepoLayout().DeleteAll(base); err != nil {
		u.CancelWithError(err)
	}
}

func (c Deinit) getPermission(u *repo_local.Repo) (success bool) {
	var err error
	ui.Err().Printf(
		"are you sure you want to deinit in %q? (y/*)",
		u.GetRepoLayout().Dir(),
	)

	var answer rune
	var n int

	if n, err = fmt.Scanf("%c", &answer); err != nil {
		ui.Err().Printf("failed to read answer: %s", err)
		return
	}

	if n != 1 {
		ui.Err().Printf("failed to read at exactly 1 answer: %s", err)
		return
	}

	if answer == 'y' || answer == 'Y' {
		success = true
	}

	return
}
