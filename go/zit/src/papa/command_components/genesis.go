package command_components

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/repo_layout"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
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
) *read_write_repo_local.Repo {
	local := read_write_repo_local.Genesis(
		c.BigBang,
		context,
		config,
		envOptions,
	)

	return local
}
