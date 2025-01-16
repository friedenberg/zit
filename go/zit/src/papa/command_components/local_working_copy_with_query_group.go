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
	req command.Request,
	builderOptions query.BuilderOptions,
) (*local_working_copy.Repo, *query.Group) {
	localWorkingCopy := cmd.MakeLocalWorkingCopy(req)

	queryGroup := cmd.MakeQueryGroup(
		req,
		builderOptions,
		localWorkingCopy,
		req.Args(),
	)

	return localWorkingCopy, queryGroup
}
