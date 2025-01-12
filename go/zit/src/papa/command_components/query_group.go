package command_components

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
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
	command any,
	repo repo.ReadWrite,
	args ...string,
) (qg *query.Group) {
	var err error

	if qg, err = repo.MakeQueryGroup(
		command,
		c.RepoId,
		c.ExternalQueryOptions,
		args...,
	); err != nil {
		repo.GetEnv().CancelWithError(err)
	}

	return
}
