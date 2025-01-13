package commands

import (
	"flag"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/inventory_list_blobs"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
	"code.linenisgreat.com/zit/go/zit/src/november/read_write_repo_local"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

// Switch to External store
type Import struct {
	immutable_config.StoreVersion
	InventoryList string
	command_components.RemoteBlobStore
	PrintCopies bool
	sku.Proto
}

func init() {
	registerCommand(
		"import",
		func(f *flag.FlagSet) CommandWithRepo {
			c := &Import{
				StoreVersion: immutable_config.CurrentStoreVersion,
			}

			f.Var(&c.StoreVersion, "store-version", "")
			f.StringVar(&c.InventoryList, "inventory-list", "", "")
			c.RemoteBlobStore.SetFlagSet(f)
			f.BoolVar(&c.PrintCopies, "print-copies", true, "output when blobs are copied")

			c.Proto.SetFlagSet(f)

			return c
		},
	)
}

func (c Import) RunWithRepo(local *read_write_repo_local.Repo, args ...string) {
	if c.InventoryList == "" {
		local.CancelWithBadRequestf("empty inventory list")
		return
	}

	bf := local.GetStore().GetInventoryListStore().FormatForVersion(c.StoreVersion)

	var rc io.ReadCloser

	// setup inventory list reader
	{
		o := dir_layout.FileReadOptions{
			Config: dir_layout.MakeConfig(
				c.Config.GetAgeEncryption(),
				c.Config.GetCompressionType(),
				false,
			),
			Path: c.InventoryList,
		}

		var err error

		if rc, err = dir_layout.NewFileReader(o); err != nil {
			local.CancelWithError(err)
		}

		defer local.MustClose(rc)
	}

	list := sku.MakeList()

	// TODO determine why this is not erroring for invalid input
	if err := inventory_list_blobs.ReadInventoryListBlob(
		bf,
		rc,
		list,
	); err != nil {
		local.CancelWithError(err)
	}

	importerOptions := store.ImporterOptions{
		CheckedOutPrinter: local.PrinterCheckedOutConflictsForRemoteTransfers(),
	}

	if c.Blobs != "" {
		{
			var err error

			if importerOptions.RemoteBlobStore, err = c.MakeRemoteBlobStore(
				local.Env,
			); err != nil {
				local.CancelWithError(err)
			}
		}
	}

	importerOptions.PrintCopies = c.PrintCopies
	importer := local.MakeImporter(importerOptions)

	if err := local.ImportList(
		list,
		importer,
	); err != nil {
		if !errors.Is(err, store.ErrNeedsMerge) {
			err = errors.Wrap(err)
		}

		local.CancelWithError(err)
	}
}
