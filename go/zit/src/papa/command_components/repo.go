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

type Repo struct{}

func (cmd *Repo) SetFlagSet(f *flag.FlagSet) {
}

func (c Repo) MakeLocalWorkingCopy(
	context *errors.Context,
	config config_mutable_cli.Config,
	envOptions env.Options,
	repoOptions local_working_copy.Options,
) *local_working_copy.Repo {
	layout := dir_layout.MakeDefault(
		context,
		config.Debug,
	)

	env := env.Make(
		context,
		config,
		layout,
		envOptions,
	)

	return local_working_copy.Make(env, repoOptions)
}

func (c Repo) MakeLocalArchive(
	repoLayout repo_layout.Layout,
) *inventory_list_store.Store {
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
		repoLayout.CancelWithError(err)
	}

	return &inventoryListStore
}

func (c Repo) MakeArchive(
	context *errors.Context,
	config config_mutable_cli.Config,
	envOptions env.Options,
	repoOptions local_working_copy.Options,
) repo.Archive {
	layout := dir_layout.MakeDefault(
		context,
		config.Debug,
	)

	env := env.Make(
		context,
		config,
		layout,
		envOptions,
	)

	layoutOptions := repo_layout.Options{
		BasePath: env.GetCLIConfig().BasePath,
	}

	var repoLayout repo_layout.Layout

	{
		var err error

		if repoLayout, err = repo_layout.Make(
			env,
			layoutOptions,
		); err != nil {
			env.CancelWithError(err)
		}
	}

	repoType := repoLayout.GetConfig().GetRepoType()

	switch repoLayout.GetConfig().GetRepoType() {
	case repo_type.TypeArchive:
		return c.MakeLocalArchive(repoLayout)

	case repo_type.TypeWorkingCopy:
		return local_working_copy.MakeWithLayout(repoOptions, repoLayout)

	default:
		env.CancelWithBadRequestf("unsupported repo type: %q", repoType)
		return nil
	}
}
