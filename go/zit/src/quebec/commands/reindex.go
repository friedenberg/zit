package commands

import (
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	command.Register("reindex", &Reindex{})
}

type Reindex struct {
	command_components.LocalWorkingCopy
}

func (cmd Reindex) Run(dep command.Request) {
	args := dep.Args()

	if len(args) > 0 {
		dep.CancelWithErrorf("reindex does not support arguments")
	}

	localWorkingCopy := cmd.MakeLocalWorkingCopyWithOptions(
		dep,
		env.Options{},
		local_working_copy.OptionsAllowConfigReadError,
	)

	localWorkingCopy.Must(localWorkingCopy.Lock)
	localWorkingCopy.Must(localWorkingCopy.GetConfig().Reset)
	localWorkingCopy.Must(localWorkingCopy.GetStore().Reindex)
	localWorkingCopy.Must(localWorkingCopy.Unlock)
}
