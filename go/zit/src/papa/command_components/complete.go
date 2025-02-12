package command_components

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type Complete struct {
	QueryGroup
}

func (cmd Complete) CompleteObjects(
	req command.Request,
	local *local_working_copy.Repo,
	queryBuilderOptions query.BuilderOption,
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
