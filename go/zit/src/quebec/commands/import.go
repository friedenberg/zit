package commands

import (
	"flag"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/box_format"
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
	hasConflicts := false

	if c.InventoryList == "" {
		err = errors.Errorf("empty inventory list")
		return
	}

	var ag age.Age

	if err = ag.AddIdentity(c.AgeIdentity); err != nil {
		err = errors.Wrapf(err, "age-identity: %q", &c.AgeIdentity)
		return
	}

	coPrinter := u.PrinterCheckedOut(box_format.CheckedOutHeaderState{})

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

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	importer := store.Importer{
		Store:      u.GetStore(),
		ErrPrinter: coPrinter,
	}

	if c.PrintCopies {
		importer.BlobCopierDelegate = func(result store.BlobCopyResult) error {
			// TODO switch to Err and fix test
			return ui.Out().Printf(
				"copied Blob %s (%d bytes)",
				result.GetBlobSha(),
				result.N,
			)
		}
	}

	if c.Blobs != "" {
		importer.RemoteBlobStore = dir_layout.MakeBlobStore(
			c.Blobs,
			&ag,
			c.CompressionType,
		)
	}

	var co *sku.CheckedOut

	for {
		sk, ok := list.Pop()

		if !ok {
			break
		}

		if co, err = importer.Import(
			sk,
		); err != nil {
			err = errors.Wrapf(err, "Sku: %s", sk)
			return
		}

		if co.GetState() == checked_out_state.Conflicted {
			hasConflicts = true

			if err = coPrinter(co); err != nil {
				err = errors.Wrap(err)
				return
			}

			continue
		}
	}

	if err = u.Unlock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if hasConflicts {
		err = store.ErrNeedsMerge
	}

	return
}
