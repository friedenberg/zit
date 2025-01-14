package commands

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

type commandWithQuery struct {
	WithQuery
	command_components.QueryGroup
}

func (c commandWithQuery) Complete(
	local *local_working_copy.Repo,
	args ...string,
) {
	if _, ok := c.WithQuery.(command_components.CompletionGenresGetter); !ok {
		return
	}

	w := sku_fmt.MakeWriterComplete(os.Stdout)
	defer local.MustClose(w)

	qg := c.MakeQueryGroup(
		query.MakeBuilderOptions(c.WithQuery),
		local,
		args...,
	)

	if err := local.GetStore().QueryTransacted(
		qg,
		w.WriteOneTransacted,
	); err != nil {
		local.CancelWithError(err)
	}
}

func (c commandWithQuery) RunWithLocalWorkingCopy(
	local *local_working_copy.Repo,
	args ...string,
) {
	qg := c.MakeQueryGroup(
		query.MakeBuilderOptions(c.WithQuery),
		local,
		args...,
	)

	defer local.PrintMatchedDormantIfNecessary()

	c.RunWithQuery(local, qg)
}
