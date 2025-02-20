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

func (cmd Genesis) OnTheFirstDay(
	req command.Request,
	repoId string,
) repo.LocalRepo {
	ui := env_ui.Make(
		req,
		req.Config,
		env_ui.Options{},
	)

	if err := cmd.Config.RepoId.Set(repoId); err != nil {
		ui.CancelWithError(err)
	}

	dir := env_dir.MakeDefaultAndInitialize(
		req,
		req.Config.Debug,
		cmd.OverrideXDGWithCwd,
	)

	var envRepo env_repo.Env

	layoutOptions := env_repo.Options{
		BasePath:             req.Config.BasePath,
		PermitNoZitDirectory: true,
	}

	{
		var err error

		if envRepo, err = env_repo.Make(
			env_local.Make(ui, dir),
			layoutOptions,
		); err != nil {
			ui.CancelWithError(err)
		}
	}

	envRepo.Genesis(cmd.BigBang)

	switch cmd.BigBang.Config.RepoType {
	case repo_type.TypeWorkingCopy:
		return local_working_copy.Genesis(
			cmd.BigBang,
			envRepo,
		)

	case repo_type.TypeArchive:
		return cmd.MakeLocalArchive(envRepo)

	default:
		req.CancelWithError(
			repo_type.ErrUnsupportedRepoType{Actual: cmd.BigBang.Config.RepoType},
		)
	}

	return nil
}
