package command_components

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type LocalWorkingCopyWithQueryGroup struct {
	LocalWorkingCopy
	QueryGroup
}

func (cmd *LocalWorkingCopyWithQueryGroup) SetFlagSet(f *flag.FlagSet) {
	cmd.LocalWorkingCopy.SetFlagSet(f)
	cmd.QueryGroup.SetFlagSet(f)
}

func (cmd LocalWorkingCopyWithQueryGroup) MakeLocalWorkingCopyAndQueryGroup(
	dep command.Dep,
	builderOptions query.BuilderOptions,
) (*local_working_copy.Repo, *query.Group) {
	localWorkingCopy := cmd.MakeLocalWorkingCopy(dep)

	queryGroup := cmd.MakeQueryGroup(
		builderOptions,
		localWorkingCopy,
		dep.Args(),
	)

	return localWorkingCopy, queryGroup
}
