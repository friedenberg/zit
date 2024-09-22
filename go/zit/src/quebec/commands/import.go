package commands

import (
	"flag"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/id"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/inventory_list"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

// Switch to External store
type Import struct {
	InventoryList   string
	Blobs           string
	AgeIdentity     age.Identity
	CompressionType immutable_config.CompressionType
	sku.Proto
}

func init() {
	registerCommand(
		"import",
		func(f *flag.FlagSet) Command {
			c := &Import{
				CompressionType: immutable_config.CompressionTypeDefault,
			}

			f.StringVar(&c.InventoryList, "inventory-list", "", "")
			f.StringVar(&c.Blobs, "blobs", "", "")
			f.Var(&c.AgeIdentity, "age-identity", "")
			c.CompressionType.AddToFlagSet(f)

			c.Proto.AddToFlagSet(f)

			return c
		},
	)
}

func (c Import) Run(u *env.Env, args ...string) (err error) {
	hasConflicts := false

	if c.InventoryList == "" {
		err = errors.Errorf("empty inventory list")
		return
	}

	if c.Blobs == "" {
		err = errors.Errorf("empty blob store")
		return
	}

	var ag age.Age

	if err = ag.AddIdentity(c.AgeIdentity); err != nil {
		err = errors.Wrapf(err, "age-identity: %q", &c.AgeIdentity)
		return
	}

	coPrinter := u.PrinterCheckedOut()

	ofo := object_inventory_format.Options{Tai: true, Verzeichnisse: true}

	bf := inventory_list.MakeFormat(u.GetConfig().GetStoreVersion(), ofo)

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

	list := inventory_list.MakeInventoryList()

	if _, err = bf.ParseBlob(rc, list); err != nil {
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
			err = errors.Wrapf(err, "Sku: %s, %#v", sk, err)
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
	blobSha := co.External.GetBlobSha()

	if u.GetFSHome().HasBlob(u.GetConfig().GetStoreVersion(), blobSha) {
		return
	}

	p := id.Path(blobSha, c.Blobs)

	o := fs_home.FileReadOptions{
		Age:             ag,
		Path:            p,
		CompressionType: c.CompressionType,
	}

	var rc sha.ReadCloser

	if rc, err = fs_home.NewFileReader(o); err != nil {
		if errors.IsNotExist(err) {
			co.SetError(errors.New("blob missing"))
			err = coErrPrinter(co)
		} else {
			err = errors.Wrapf(err, "Path: %q", p)
		}

		return
	}

	defer errors.DeferredCloser(&err, rc)

	var aw sha.WriteCloser

	if aw, err = u.GetFSHome().BlobWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	var n int64

	if n, err = io.Copy(aw, rc); err != nil {
		co.SetError(errors.New("blob copy failed"))
		err = coErrPrinter(co)
		return
	}

	shaRc := rc.GetShaLike()

	if !shaRc.EqualsSha(blobSha) {
		co.SetError(errors.New("blob sha mismatch"))
		err = coErrPrinter(co)
		ui.TodoRecoverable(
			"sku blob mismatch: sku had %s while blob store had %s",
			co.Internal.GetBlobSha(),
			shaRc,
		)
	}

	ui.Err().Printf("copied Blob %s (%d bytes)", blobSha, n)

	return
}
