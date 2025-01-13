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
	"code.linenisgreat.com/zit/go/zit/src/november/read_write_repo_local"
)

type Genesis struct {
	repo_layout.BigBang
}

func (cmd *Genesis) SetFlagSet(f *flag.FlagSet) {
	cmd.BigBang.SetFlagSet(f)
}

func (c Genesis) OnTheFirstDay(
	context *errors.Context,
	config config_mutable_cli.Config,
	envOptions env.Options,
) repo.ReadWrite {
	switch c.BigBang.Config.RepoType {

	case repo_type.TypeReadWrite:
		return c.readWrite(context, config, envOptions)

	case repo_type.TypeRelay:
		fallthrough
		// return c.relay(context, config, envOptions)

	default:
		context.CancelWithError(
			repo_type.ErrUnsupportedRepoType{Actual: c.BigBang.Config.RepoType},
		)
	}

	return nil
}

func (c Genesis) readWrite(
	context *errors.Context,
	config config_mutable_cli.Config,
	envOptions env.Options,
) repo.ReadWrite {
	local := read_write_repo_local.Genesis(
		c.BigBang,
		context,
		config,
		envOptions,
	)

	return local
}

func (c Genesis) relay(
	context *errors.Context,
	config config_mutable_cli.Config,
	envOptions env.Options,
) repo.Relay {
	layout := dir_layout.MakeDefault(
		context,
		config.Debug,
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

	return repoLayout
}
