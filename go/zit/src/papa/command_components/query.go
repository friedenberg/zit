package command_components

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	pkg_query "code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
)

type Query struct {
	sku.ExternalQueryOptions
}

func (cmd *Query) SetFlagSet(f *flag.FlagSet) {
	// TODO switch to repo
	f.Var(&cmd.RepoId, "kasten", "none or Browser")
	f.BoolVar(&cmd.ExcludeUntracked, "exclude-untracked", false, "")
	f.BoolVar(&cmd.ExcludeRecognized, "exclude-recognized", false, "")
}

func (cmd Query) MakeQuery(
	req command.Request,
	options pkg_query.BuilderOption,
	workingCopy repo.WorkingCopy,
	args []string,
) (query *pkg_query.Query) {
	if repo, ok := workingCopy.(repo.LocalWorkingCopy); ok {
		envWorkspace := repo.GetEnvWorkspace()

		options = pkg_query.BuilderOptions(
			options,
			pkg_query.BuilderOptionWorkspace{Env: envWorkspace},
		)
	}

	var err error

	if query, err = workingCopy.MakeExternalQueryGroup(
		options,
		cmd.ExternalQueryOptions,
		args...,
	); err != nil {
		req.CancelWithError(err)
	}

	return
}
