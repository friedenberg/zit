package commands

import (
	"flag"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/inventory_list_blobs"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	command.Register(
		"import",
		&Import{
			StoreVersion: immutable_config.CurrentStoreVersion,
		},
	)
}

// Switch to External store
type Import struct {
	command_components.LocalWorkingCopy
	command_components.RemoteBlobStore

	immutable_config.StoreVersion
	InventoryList string
	PrintCopies   bool
	sku.Proto
}

func (cmd *Import) SetFlagSet(f *flag.FlagSet) {
	f.Var(&cmd.StoreVersion, "store-version", "")
	f.StringVar(&cmd.InventoryList, "inventory-list", "", "")
	cmd.RemoteBlobStore.SetFlagSet(f)
	f.BoolVar(&cmd.PrintCopies, "print-copies", true, "output when blobs are copied")

	cmd.Proto.SetFlagSet(f)
}

func (cmd Import) Run(dep command.Dep) {
	localWorkingCopy := cmd.MakeLocalWorkingCopy(dep)

	if cmd.InventoryList == "" {
		dep.CancelWithBadRequestf("empty inventory list")
	}

	bf := localWorkingCopy.GetStore().GetInventoryListStore().FormatForVersion(cmd.StoreVersion)

	var rc io.ReadCloser

	// setup inventory list reader
	{
		o := dir_layout.FileReadOptions{
			Config: dir_layout.MakeConfig(
				cmd.Config.GetAgeEncryption(),
				cmd.Config.GetCompressionType(),
				false,
			),
			Path: cmd.InventoryList,
		}

		var err error

		if rc, err = dir_layout.NewFileReader(o); err != nil {
			localWorkingCopy.CancelWithError(err)
		}

		defer localWorkingCopy.MustClose(rc)
	}

	list := sku.MakeList()

	// TODO determine why this is not erroring for invalid input
	if err := inventory_list_blobs.ReadInventoryListBlob(
		bf,
		rc,
		list,
	); err != nil {
		localWorkingCopy.CancelWithError(err)
	}

	importerOptions := store.ImporterOptions{
		CheckedOutPrinter: localWorkingCopy.PrinterCheckedOutConflictsForRemoteTransfers(),
	}

	if cmd.Blobs != "" {
		{
			var err error

			if importerOptions.RemoteBlobStore, err = cmd.MakeRemoteBlobStore(
				localWorkingCopy.Env,
			); err != nil {
				localWorkingCopy.CancelWithError(err)
			}
		}
	}

	importerOptions.PrintCopies = cmd.PrintCopies
	importer := localWorkingCopy.MakeImporter(importerOptions)

	if err := localWorkingCopy.ImportList(
		list,
		importer,
	); err != nil {
		if !errors.Is(err, store.ErrNeedsMerge) {
			err = errors.Wrap(err)
		}

		localWorkingCopy.CancelWithError(err)
	}
}
