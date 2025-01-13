package command_components

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/charlie/options_print"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/hotel/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/repo_layout"
	"code.linenisgreat.com/zit/go/zit/src/kilo/box_format"
	"code.linenisgreat.com/zit/go/zit/src/lima/blob_store"
	"code.linenisgreat.com/zit/go/zit/src/lima/inventory_list_store"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
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
) repo.Archive {
	switch c.BigBang.Config.RepoType {
	case repo_type.TypeWorkingCopy:
		return c.makeWorkingCopy(context, config, envOptions)

	case repo_type.TypeArchive:
		return c.makeArchive(context, config, envOptions)

	default:
		context.CancelWithError(
			repo_type.ErrUnsupportedRepoType{Actual: c.BigBang.Config.RepoType},
		)
	}

	return nil
}

func (c Genesis) makeWorkingCopy(
	context *errors.Context,
	config config_mutable_cli.Config,
	envOptions env.Options,
) repo.WorkingCopy {
	local := local_working_copy.Genesis(
		c.BigBang,
		context,
		config,
		envOptions,
	)

	return local
}

func (c Genesis) makeArchive(
	context *errors.Context,
	config config_mutable_cli.Config,
	envOptions env.Options,
) repo.Archive {
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

	objectFormat := object_inventory_format.FormatForVersion(repoLayout.GetStoreVersion())
	boxFormat := box_format.MakeBoxTransactedArchive(
		repoLayout.Env,
		options_print.V0{}.WithPrintTai(true),
	)

	inventoryListBlobStore := blob_store.MakeInventoryStore(
		repoLayout,
		objectFormat,
		boxFormat,
	)

	var inventoryListStore inventory_list_store.Store

	if err := inventoryListStore.Initialize(
		repoLayout,
		objectFormat,
		nil,
		inventoryListBlobStore,
	); err != nil {
		env.CancelWithError(err)
	}

	return &inventoryListStore
}
