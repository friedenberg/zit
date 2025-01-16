package commands

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	command.Register("find-missing", &FindMissing{})
}

type FindMissing struct {
	command_components.LocalWorkingCopy
}

func (cmd FindMissing) Run(dep command.Request) {
	localWorkingCopy := cmd.MakeLocalWorkingCopy(dep)

	var lookupStored map[sha.Bytes][]string

	{
		var err error

		if lookupStored, err = localWorkingCopy.GetStore().MakeBlobShaBytesMap(); err != nil {
			dep.CancelWithError(err)
		}
	}

	for _, shSt := range dep.Args() {
		var sh sha.Sha

		if err := sh.Set(shSt); err != nil {
			localWorkingCopy.CancelWithError(err)
		}

		oids, ok := lookupStored[sh.GetBytes()]

		if ok {
			localWorkingCopy.GetUI().Printf("%s (checked in as %q)", &sh, oids)
		} else {
			localWorkingCopy.GetUI().Printf("%s (missing)", &sh)
		}
	}
}
