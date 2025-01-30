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
	command.Register("deinit", &Deinit{})
}

type Deinit struct {
	command_components.LocalWorkingCopy

	Force bool
}

func (cmd *Deinit) SetFlagSet(f *flag.FlagSet) {
	f.BoolVar(
		&cmd.Force,
		"force",
		false,
		"force deinit",
	)
}

func (cmd Deinit) Run(dep command.Request) {
	// TODO switch to archive
	localWorkingCopy := cmd.MakeLocalWorkingCopy(dep)

	if !cmd.Force && !cmd.getPermission(localWorkingCopy) {
		return
	}

	base := path.Join(localWorkingCopy.GetEnvRepo().Dir())

	if err := files.SetAllowUserChangesRecursive(base); err != nil {
		localWorkingCopy.CancelWithError(err)
	}

	if err := localWorkingCopy.GetEnvRepo().DeleteAll(base); err != nil {
		localWorkingCopy.CancelWithError(err)
	}
}

func (c Deinit) getPermission(u *local_working_copy.Repo) (success bool) {
	var err error
	ui.Err().Printf(
		"are you sure you want to deinit in %q? (y/*)",
		u.GetEnvRepo().Dir(),
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
