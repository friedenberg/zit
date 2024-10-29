package commands

import (
	"flag"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_metadata_fmt"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/inventory_list_fmt"
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

			c.Proto.AddToFlagSet(f)

			return c
		},
	)
}

func (c Import) Run(u *env.Env, args ...string) (result Result) {
	result.Error = c.run(u, args...)

	return
}

func (c Import) run(u *env.Env, args ...string) (err error) {
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

	coPrinter := u.PrinterCheckedOut()

	bf := u.GetStore().GetInventoryListStore().FormatForVersion(c.StoreVersion)

	var rc io.ReadCloser

	// setup inventory list reader
	{
		o := fs_home.FileReadOptions{
			Age:             &ag,
			Path:            c.InventoryList,
			CompressionType: c.CompressionType,
		}

		if rc, err = fs_home.NewFileReader(o); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, rc)
	}

	list := sku.MakeList()

	// TODO determine why this is not erroring for invalid input
	if err = inventory_list_fmt.ReadInventoryListBlob(bf, rc, list); err != nil {
		err = errors.Wrap(err)
		return
	}

	u.Lock()
	defer u.Unlock()

	var co *sku.CheckedOut

	for {
		sk, ok := list.Pop()

		if !ok {
			break
		}

		if co, err = u.GetStore().Import(
			sk,
		); err != nil {
			err = errors.Wrapf(err, "Sku: %s", sk)
			return
		}

		if co.State == checked_out_state.Conflicted {
			hasConflicts = true

			if err = coPrinter(co); err != nil {
				err = errors.Wrap(err)
				return
			}

			continue
		}

		if err = c.importBlobIfNecessary(u, co, &ag, coPrinter); err != nil {
			if age.IsNoIdentityMatchError(err) {
				err = nil
			} else {
				err = errors.Wrapf(err, "Checked Out: %q", co)
				return
			}
		}

		if co.State == checked_out_state.Error {
			co.External.Metadata.Fields = append(
				co.External.Metadata.Fields,
				object_metadata_fmt.MetadataFieldError(co.Error)...,
			)

			if err = coPrinter(co); err != nil {
				err = errors.Wrap(err)
				return
			}

			continue
		}
	}

	if hasConflicts {
		err = store.ErrNeedsMerge
	}

	return
}

func (c Import) importBlobIfNecessary(
	u *env.Env,
	co *sku.CheckedOut,
	ag *age.Age,
	coErrPrinter interfaces.FuncIter[*sku.CheckedOut],
) (err error) {
	if c.Blobs == "" {
		return
	}

	blobStore := fs_home.MakeBlobStore(
		c.Blobs,
		ag,
		c.CompressionType,
	)

	blobSha := co.External.GetBlobSha()

	var n int64

	if n, err = u.GetFSHome().CopyBlobIfNecessary(blobStore, blobSha); err != nil {
		if errors.Is(err, &fs_home.ErrAlreadyExists{}) {
			err = nil
		} else {
			co.SetError(err)
			err = coErrPrinter(co)
		}
		return
	}

	// TODO switch to Err and fix test
	ui.Out().Printf("copied Blob %s (%d bytes)", blobSha, n)

	return
}
