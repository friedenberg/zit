package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	command.Register(
		"init-workspace",
		&InitWorkspace{},
	)
}

type InitWorkspace struct {
	command_components.LocalWorkingCopy
}

func (cmd *InitWorkspace) SetFlagSet(f *flag.FlagSet) {
	cmd.LocalWorkingCopy.SetFlagSet(f)
}

func (cmd InitWorkspace) Run(req command.Request) {
	localWorkingCopy := cmd.MakeLocalWorkingCopy(req)

	if err := localWorkingCopy.GetEnvWorkspace().CreateWorkspace(); err != nil {
		req.CancelWithError(err)
	}
}
