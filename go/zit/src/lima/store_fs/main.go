package store_fs

import (
	"encoding/gob"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/box"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
	"code.linenisgreat.com/zit/go/zit/src/hotel/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/env_workspace"
)

func init() {
	gob.Register(sku.Transacted{})
}

func Make(
	config sku.Config,
	deletedPrinter interfaces.FuncIter[*fd.FD],
	fileExtensions interfaces.FileExtensionGetter,
	envRepo env_repo.Env,
	inventoryFormatOptions object_inventory_format.Options,
	fileEncoder FileEncoder,
	envWorkspace env_workspace.Env,
) (fs *Store, err error) {
	fs = &Store{
		config:         config,
		deletedPrinter: deletedPrinter,
		envRepo:        envRepo,
		envWorkspace:   envWorkspace,
		fileEncoder:    fileEncoder,
		fileExtensions: fileExtensions,
		dir:            envRepo.GetCwd(),
		dirItems: makeObjectsWithDir(
			envWorkspace.GetWorkspaceDir(),
			fileExtensions,
			envRepo,
		),
		deleted: collections_value.MakeMutableValueSet[*fd.FD](
			nil,
		),
		deletedInternal: collections_value.MakeMutableValueSet[*fd.FD](
			nil,
		),
		objectFormatOptions: inventoryFormatOptions,
		metadataTextParser: object_metadata.MakeTextParser(
			object_metadata.Dependencies{
				EnvDir:    envRepo,
				BlobStore: envRepo,
			},
		),
	}

	return
}

type Store struct {
	config              sku.Config
	deletedPrinter      interfaces.FuncIter[*fd.FD]
	metadataTextParser  object_metadata.TextParser
	envRepo             env_repo.Env
	envWorkspace        env_workspace.Env
	fileEncoder         FileEncoder
	inlineTypeChecker   ids.InlineTypeChecker
	fileExtensions      interfaces.FileExtensionGetter
	dir                 string
	objectFormatOptions object_inventory_format.Options

	dirItems

	deleteLock      sync.Mutex
	deleted         fd.MutableSet
	deletedInternal fd.MutableSet
}

func (fs *Store) GetExternalStoreLike() env_workspace.StoreLike {
	return fs
}

// Deletions of user objects that should be exposed to the user
func (s *Store) DeleteCheckedOut(co *sku.CheckedOut) (err error) {
	external := co.GetSkuExternal()

	var item *sku.FSItem

	if item, err = s.ReadFSItemFromExternal(external); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.deleteLock.Lock()
	defer s.deleteLock.Unlock()

	if err = item.MutableSetLike.Each(s.deleted.Add); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// Deletions of "transient" internal objects that should not be exposed to the
// user
func (s *Store) DeleteCheckedOutInternal(co *sku.CheckedOut) (err error) {
	external := co.GetSkuExternal()

	var i *sku.FSItem

	if i, err = s.ReadFSItemFromExternal(external); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.deleteLock.Lock()
	defer s.deleteLock.Unlock()

	if err = i.MutableSetLike.Each(s.deletedInternal.Add); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (fs *Store) Flush() (err error) {
	deleteOp := DeleteCheckout{}

	if err = deleteOp.Run(
		fs.config.IsDryRun(),
		fs.envRepo,
		fs.deletedPrinter,
		fs.deleted,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = deleteOp.Run(
		fs.config.IsDryRun(),
		fs.envRepo,
		nil,
		fs.deletedInternal,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	fs.deleted.Reset()
	fs.deletedInternal.Reset()

	return
}

func (fs *Store) String() (out string) {
	if quiter.Len(
		fs.dirItems.probablyCheckedOut,
		fs.definitelyNotCheckedOut,
	) == 0 {
		return
	}

	sb := &strings.Builder{}
	sb.WriteRune(box.OpGroupOpen)

	hasOne := false

	writeOneIfNecessary := func(v interfaces.Stringer) (err error) {
		if hasOne {
			sb.WriteRune(box.OpOr)
		}

		sb.WriteString(v.String())

		hasOne = true

		return
	}

	fs.dirItems.probablyCheckedOut.Each(
		func(z *sku.FSItem) (err error) {
			return writeOneIfNecessary(z)
		},
	)

	fs.definitelyNotCheckedOut.Each(
		func(z *sku.FSItem) (err error) {
			return writeOneIfNecessary(z)
		},
	)

	sb.WriteRune(box.OpGroupClose)

	out = sb.String()
	return
}

func (s *Store) GetExternalObjectIds() (ks []*sku.FSItem, err error) {
	if err = s.dirItems.processRootDir(); err != nil {
		err = errors.Wrap(err)
		return
	}

	ks = make([]*sku.FSItem, 0)
	var l sync.Mutex

	if err = s.All(
		func(kfp *sku.FSItem) (err error) {
			l.Lock()
			defer l.Unlock()

			ks = append(ks, kfp)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) GetFSItemsForDir(
	fd *fd.FD,
) (items []*sku.FSItem, err error) {
	if !fd.IsDir() {
		err = errors.Errorf("not a directory: %q", fd)
		return
	}

	if items, err = s.dirItems.processDir(fd.GetPath()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO confirm against actual Object Id
func (s *Store) GetFSItemsForString(
	baseDir string,
	value string,
	tryPattern bool,
) (items []*sku.FSItem, err error) {
	if value == "." {
		if items, err = s.GetExternalObjectIds(); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	var fdee *fd.FD

	if fdee, err = fd.MakeFromPath(baseDir, value, s.envRepo); err != nil {
		if errors.IsNotExist(err) && tryPattern {
			if items, err = s.dirItems.processFDPattern(
				value,
				filepath.Join(s.dir, fmt.Sprintf("%s*", value)),
				s.dir,
			); err != nil {
				err = errors.Wrap(err)
				return
			}
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	if fdee.IsDir() {
		if items, err = s.GetFSItemsForDir(fdee); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if _, items, err = s.dirItems.processFD(fdee); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (store *Store) GetObjectIdsForString(
	v string,
) (k []sku.ExternalObjectId, err error) {
	var items []*sku.FSItem

	if items, err = store.GetFSItemsForString(store.root, v, false); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, item := range items {
		var eoid ids.ExternalObjectId

		if err = item.WriteToExternalObjectId(
			&eoid,
			store.envRepo,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		k = append(k, &eoid)
	}

	return
}

func (fs *Store) Get(
	k interfaces.ObjectId,
) (t *sku.FSItem, ok bool) {
	return fs.dirItems.probablyCheckedOut.Get(k.String())
}

func (s *Store) Initialize(esi env_workspace.Supplies) (err error) {
	s.externalStoreSupplies = esi
	return
}

func (s *Store) ReadAllExternalItems() (err error) {
	if err = s.dirItems.processRootDir(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadFSItemFromExternal(
	tg sku.TransactedGetter,
) (item *sku.FSItem, err error) {
	item = &sku.FSItem{} // TODO use pool or use dir_items?
	item.Reset()

	sk := tg.GetSku()

	// TODO handle sort order
	for _, field := range sk.Metadata.Fields {
		var fdee *fd.FD

		switch strings.ToLower(field.Key) {
		case "object":
			fdee = &item.Object

		case "blob":
			fdee = &item.Blob

		case "conflict":
			fdee = &item.Conflict

		default:
			err = errors.Errorf("unexpected field: %#v", field)
			return
		}

		// if we've already set one of object, blob, or conflict, don't set it again
		// and instead add a new FD to the item
		if !fdee.IsEmpty() {
			fdee = &fd.FD{}
		}

		if err = fdee.SetIgnoreNotExists(field.Value); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = item.Add(fdee); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = item.ExternalObjectId.SetObjectIdLike(
		&sk.ObjectId,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	// external.ObjectId.ResetWith(conflicted.GetSkuExternal().GetObjectId())
	// TODO populate FD
	if !sk.ExternalObjectId.IsEmpty() {
		if err = item.ExternalObjectId.SetObjectIdLike(
			&item.ExternalObjectId,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) WriteFSItemToExternal(
	item *sku.FSItem,
	tg sku.TransactedGetter,
) (err error) {
	external := tg.GetSku()
	external.Metadata.Fields = external.Metadata.Fields[:0]

	m := &external.Metadata
	m.Tai = item.GetTai()

	var mode checkout_mode.Mode

	if mode, err = item.GetCheckoutModeOrError(); err != nil {
		if sku.IsErrMergeConflict(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	switch mode {
	case checkout_mode.BlobOnly:
		before := item.Blob.String()
		after := s.envRepo.Rel(before)

		if err = external.ExternalObjectId.SetBlob(after); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		k := &item.ExternalObjectId

		if err = external.ObjectId.SetObjectIdLike(k); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = external.ExternalObjectId.SetObjectIdLike(
			&item.ExternalObjectId,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if external.ExternalObjectId.String() != k.String() {
			err = errors.Errorf(
				"expected %q but got %q. %s",
				k,
				&external.ExternalObjectId,
				item.Debug(),
			)

			return
		}
	}

	if err = item.WriteToSku(
		external,
		s.envRepo,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	fdees := quiter.SortedValues(item.MutableSetLike)

	for _, f := range fdees {
		field := object_metadata.Field{
			Value:     f.GetPath(),
			ColorType: string_format_writer.ColorTypeId,
		}

		switch {
		case item.Object.Equals(f):
			field.Key = "object"

		case item.Conflict.Equals(f):
			field.Key = "conflict"

		case item.Blob.Equals(f):
			fallthrough

		default:
			field.Key = "blob"
		}

		external.Metadata.Fields = append(external.Metadata.Fields, field)
	}

	return
}
