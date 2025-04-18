package command_components

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type LocalWorkingCopyWithQueryGroup struct {
	LocalWorkingCopy
	Query
}

func (cmd *LocalWorkingCopyWithQueryGroup) SetFlagSet(f *flag.FlagSet) {
	cmd.LocalWorkingCopy.SetFlagSet(f)
	cmd.Query.SetFlagSet(f)
}

func (cmd LocalWorkingCopyWithQueryGroup) MakeLocalWorkingCopyAndQueryGroup(
	req command.Request,
	builderOptions query.BuilderOption,
) (*local_working_copy.Repo, *query.Query) {
	localWorkingCopy := cmd.MakeLocalWorkingCopy(req)

	queryGroup := cmd.MakeQueryIncludingWorkspace(
		req,
		builderOptions,
		localWorkingCopy,
		req.PopArgs(),
	)

	return localWorkingCopy, queryGroup
}
