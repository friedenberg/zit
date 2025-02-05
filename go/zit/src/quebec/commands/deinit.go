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
		ui.Err().Print("permission denied and -force not specified, aborting")
		return
	}

	base := path.Join(localWorkingCopy.GetEnvRepo().Dir())

	if err := files.SetAllowUserChangesRecursive(base); err != nil {
		localWorkingCopy.CancelWithError(err)
	}

	if err := localWorkingCopy.GetEnvRepo().Delete(
		localWorkingCopy.GetEnvRepo().GetXDG().GetXDGPaths()...,
	); err != nil {
		localWorkingCopy.CancelWithError(err)
	}

	if err := localWorkingCopy.GetEnvWorkspace().DeleteWorkspace(); err != nil {
		localWorkingCopy.CancelWithError(err)
	}
}

func (c Deinit) getPermission(repo *local_working_copy.Repo) bool {
	return repo.Confirm(
		fmt.Sprintf(
			"are you sure you want to deinit in %q?",
			repo.GetEnvRepo().Dir(),
		),
	)
}
