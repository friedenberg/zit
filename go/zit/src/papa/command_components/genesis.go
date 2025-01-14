package command_components

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/hotel/repo_layout"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type Genesis struct {
	repo_layout.BigBang
	Repo
}

func (cmd *Genesis) SetFlagSet(f *flag.FlagSet) {
	cmd.BigBang.SetFlagSet(f)
}

func (c Genesis) OnTheFirstDay(
	context *errors.Context,
	config config_mutable_cli.Config,
	envOptions env.Options,
) repo.Archive {
	layout := dir_layout.MakeDefaultAndInitialize(
		context,
		config.Debug,
		c.OverrideXDGWithCwd,
	)

	env := env.Make(
		context,
		config,
		layout,
		env.Options{},
	)

	var repoLayout repo_layout.Layout

	layoutOptions := repo_layout.Options{
		BasePath:             config.BasePath,
		PermitNoZitDirectory: true,
	}

	{
		var err error

		if repoLayout, err = repo_layout.Make(
			env,
			layoutOptions,
		); err != nil {
			env.CancelWithError(err)
		}
	}

	repoLayout.Genesis(c.BigBang)

	switch c.BigBang.Config.RepoType {
	case repo_type.TypeWorkingCopy:
		return c.makeWorkingCopy(repoLayout)

	case repo_type.TypeArchive:
		return c.makeArchive(repoLayout)

	default:
		context.CancelWithError(
			repo_type.ErrUnsupportedRepoType{Actual: c.BigBang.Config.RepoType},
		)
	}

	return nil
}

func (c Genesis) makeWorkingCopy(
	repoLayout repo_layout.Layout,
) repo.WorkingCopy {
	local := local_working_copy.Genesis(
		c.BigBang,
		repoLayout,
	)

	return local
}

func (c Genesis) makeArchive(
	repoLayout repo_layout.Layout,
) repo.Archive {
	return c.MakeLocalArchive(repoLayout)
}
