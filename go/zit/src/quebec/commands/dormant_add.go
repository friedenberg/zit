package commands

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	command.Register("dormant-add", &DormantAdd{})
}

type DormantAdd struct {
	command_components.LocalWorkingCopy
}

func (cmd DormantAdd) Run(dep command.Dep) {
	localWorkingCopy := cmd.MakeLocalWorkingCopy(dep)

	localWorkingCopy.Must(localWorkingCopy.Lock)

	for _, v := range dep.Args() {
		cs := catgut.MakeFromString(v)

		if err := localWorkingCopy.GetDormantIndex().AddDormantTag(cs); err != nil {
			localWorkingCopy.CancelWithError(err)
		}
	}

	localWorkingCopy.Must(localWorkingCopy.Unlock)
}
