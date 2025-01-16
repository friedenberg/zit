package command_components

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/hotel/repo_layout"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type Genesis struct {
	repo_layout.BigBang
	LocalWorkingCopy
	LocalArchive
}

func (cmd *Genesis) SetFlagSet(f *flag.FlagSet) {
	cmd.BigBang.SetFlagSet(f)
}

func (cmd Genesis) OnTheFirstDay(dep command.Request) repo.Archive {
	layout := dir_layout.MakeDefaultAndInitialize(
		dep.Context,
		dep.Config.Debug,
		cmd.OverrideXDGWithCwd,
	)

	env := env.Make(
		dep.Context,
		dep.Config,
		layout,
		env.Options{},
	)

	var repoLayout repo_layout.Layout

	layoutOptions := repo_layout.Options{
		BasePath:             dep.Config.BasePath,
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

	repoLayout.Genesis(cmd.BigBang)

	switch cmd.BigBang.Config.RepoType {
	case repo_type.TypeWorkingCopy:
		return local_working_copy.Genesis(
			cmd.BigBang,
			repoLayout,
		)

	case repo_type.TypeArchive:
		return cmd.MakeLocalArchive(repoLayout)

	default:
		dep.CancelWithError(
			repo_type.ErrUnsupportedRepoType{Actual: cmd.BigBang.Config.RepoType},
		)
	}

	return nil
}
