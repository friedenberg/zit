package commands

import (
	"flag"
	"fmt"
	"path"

	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	registerCommand("deinit", &Deinit{})
}

type Deinit struct {
	command_components.LocalWorkingCopy

	Force bool
}

func (c *Deinit) SetFlagSet(f *flag.FlagSet) {
	f.BoolVar(
		&c.Force,
		"force",
		false,
		"force deinit",
	)
}

func (c Deinit) Run(dep command.Dep) {
	// TODO switch to archive
	u := c.MakeLocalWorkingCopy(dep)

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

func (c Deinit) getPermission(u *local_working_copy.Repo) (success bool) {
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
