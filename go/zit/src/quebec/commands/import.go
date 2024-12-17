package commands

import (
	"flag"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/inventory_list_blobs"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
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
		func(f *flag.FlagSet) CommandWithResult {
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

func (c Import) Run(u *env.Local, args ...string) (result Result) {
	result.Error = c.run(u, args...)

	return
}

func (c Import) run(u *env.Local, args ...string) (err error) {
	if c.InventoryList == "" {
		err = errors.Errorf("empty inventory list")
		return
	}

	var ag age.Age

	if err = ag.AddIdentity(c.AgeIdentity); err != nil {
		err = errors.Wrapf(err, "age-identity: %q", &c.AgeIdentity)
		return
	}

	bf := u.GetStore().GetInventoryListStore().FormatForVersion(c.StoreVersion)

	var rc io.ReadCloser

	// setup inventory list reader
	{
		o := dir_layout.FileReadOptions{
			Age:             &ag,
			Path:            c.InventoryList,
			CompressionType: c.CompressionType,
		}

		if rc, err = dir_layout.NewFileReader(o); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, rc)
	}

	list := sku.MakeList()

	// TODO determine why this is not erroring for invalid input
	if err = inventory_list_blobs.ReadInventoryListBlob(bf, rc, list); err != nil {
		err = errors.Wrap(err)
		return
	}

	var remoteBlobStore dir_layout.BlobStore

	if c.Blobs != "" {
		remoteBlobStore = dir_layout.MakeBlobStore(
			c.Blobs,
			&ag,
			c.CompressionType,
		)
	}

	if err = u.ImportListFromRemoteBlobStore(
		list,
		remoteBlobStore,
		c.PrintCopies,
	); err != nil {
		if !errors.Is(err, store.ErrNeedsMerge) {
			err = errors.Wrap(err)
		}

		return
	}

	return
}
