package command_components

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/env_ui"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_local"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type Genesis struct {
	env_repo.BigBang
	LocalWorkingCopy
	LocalArchive
}

func (cmd *Genesis) SetFlagSet(f *flag.FlagSet) {
	cmd.BigBang.SetFlagSet(f)
}

func (cmd Genesis) OnTheFirstDay(dep command.Request) repo.Repo {
	dir := env_dir.MakeDefaultAndInitialize(
		dep,
		dep.Config.Debug,
		cmd.OverrideXDGWithCwd,
	)

	ui := env_ui.Make(
		dep,
		dep.Config,
		env_ui.Options{},
	)

	var repoLayout env_repo.Env

	layoutOptions := env_repo.Options{
		BasePath:             dep.Config.BasePath,
		PermitNoZitDirectory: true,
	}

	{
		var err error

		if repoLayout, err = env_repo.Make(
			env_local.Make(ui, dir),
			layoutOptions,
		); err != nil {
			ui.CancelWithError(err)
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
