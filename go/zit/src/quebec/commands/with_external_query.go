package commands

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

type commandWithQuery struct {
	CommandWithQuery
	command_components.QueryGroup
}

type CompletionGenresGetter interface {
	CompletionGenres() ids.Genre
}

func (c commandWithQuery) CompleteWithRepo(
	local *repo_local_working_copy.Repo,
	args ...string,
) {
	if _, ok := c.CommandWithQuery.(CompletionGenresGetter); !ok {
		return
	}

	w := sku_fmt.MakeWriterComplete(os.Stdout)
	defer local.MustClose(w)

	qg := c.MakeQueryGroup(
		query.MakeBuilderOptions(c.CommandWithQuery),
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

func (c commandWithQuery) RunWithRepo(
	local *repo_local_working_copy.Repo,
	args ...string,
) {
	qg := c.MakeQueryGroup(
		query.MakeBuilderOptions(c.CommandWithQuery),
		local,
		args...,
	)

	defer local.PrintMatchedDormantIfNecessary()

	c.RunWithQuery(local, qg)
}
