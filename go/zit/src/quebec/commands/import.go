package commands

import (
	"flag"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/repo_layout"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/inventory_list_blobs"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

// Switch to External store
type Import struct {
	immutable_config.StoreVersion
	InventoryList   string
	Blobs           string
	AgeIdentity     age.Identity
	CompressionType immutable_config.CompressionType
	PrintCopies     bool
	sku.Proto
}

func init() {
	registerCommand(
		"import",
		func(f *flag.FlagSet) CommandWithRepo {
			c := &Import{
				StoreVersion:    immutable_config.CurrentStoreVersion,
				CompressionType: immutable_config.CompressionTypeDefault,
			}

			f.Var(&c.StoreVersion, "store-version", "")
			f.StringVar(&c.InventoryList, "inventory-list", "", "")
			f.StringVar(&c.Blobs, "blobs", "", "")
			f.Var(&c.AgeIdentity, "age-identity", "")
			c.CompressionType.AddToFlagSet(f)
			f.BoolVar(&c.PrintCopies, "print-copies", true, "output when blobs are copied")

			c.Proto.AddToFlagSet(f)

			return c
		},
	)
}

func (c Import) RunWithRepo(local *repo_local.Repo, args ...string) {
	if c.InventoryList == "" {
		local.CancelWithBadRequestf("empty inventory list")
		return
	}

	var ag age.Age

	if err := ag.AddIdentity(c.AgeIdentity); err != nil {
		local.CancelWithErrorAndFormat(err, "age-identity: %q", &c.AgeIdentity)
	}

	bf := local.GetStore().GetInventoryListStore().FormatForVersion(c.StoreVersion)

	var rc io.ReadCloser

	// setup inventory list reader
	{
		o := repo_layout.FileReadOptions{
			Age:             &ag,
			Path:            c.InventoryList,
			CompressionType: c.CompressionType,
		}

		var err error

		if rc, err = repo_layout.NewFileReader(o); err != nil {
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

	importer := local.MakeImporter(c.PrintCopies)

	if c.Blobs != "" {
		importer.RemoteBlobStore = repo_layout.MakeBlobStore(
			c.Blobs,
			&ag,
			c.CompressionType,
		)
	}

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
