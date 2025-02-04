package command_components

import (
	"flag"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type QueryGroup struct {
	sku.ExternalQueryOptions
}

func (cmd *QueryGroup) SetFlagSet(f *flag.FlagSet) {
	f.Var(&cmd.RepoId, "kasten", "none or Browser")
	f.BoolVar(&cmd.ExcludeUntracked, "exclude-untracked", false, "")
	f.BoolVar(&cmd.ExcludeRecognized, "exclude-recognized", false, "")
}

func (c QueryGroup) MakeQueryGroup(
	req command.Request,
	options query.BuilderOptions,
	repo repo.WorkingCopy,
	args []string,
) (queryGroup *query.Group) {
	var err error

	if queryGroup, err = repo.MakeExternalQueryGroup(
		options,
		c.ExternalQueryOptions,
		args...,
	); err != nil {
		req.CancelWithError(err)
	}

	return
}

func (cmd QueryGroup) CompleteWithRepo(
	req command.Request,
	local *local_working_copy.Repo,
	queryBuilderOptions query.BuilderOptions,
	args ...string,
) {
	completionWriter := sku_fmt.MakeWriterComplete(os.Stdout)
	defer local.MustClose(completionWriter)

	queryGroup := cmd.MakeQueryGroup(
		req,
		queryBuilderOptions,
		local,
		args,
	)

	if err := local.GetStore().QueryTransacted(
		queryGroup,
		completionWriter.WriteOneTransacted,
	); err != nil {
		local.CancelWithError(err)
	}
}
