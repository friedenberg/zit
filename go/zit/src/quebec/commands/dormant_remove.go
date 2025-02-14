package commands

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	command.Register("dormant-remove", &DormantRemove{})
}

type DormantRemove struct {
	command_components.LocalWorkingCopy
}

func (cmd DormantRemove) Run(dep command.Request) {
	localWorkingCopy := cmd.MakeLocalWorkingCopy(dep)
	localWorkingCopy.Must(localWorkingCopy.Lock)

	for _, v := range dep.PopArgs() {
		cs := catgut.MakeFromString(v)

		if err := localWorkingCopy.GetDormantIndex().RemoveDormantTag(
			cs,
		); err != nil {
			localWorkingCopy.CancelWithError(err)
		}
	}

	localWorkingCopy.Must(localWorkingCopy.Unlock)
}
