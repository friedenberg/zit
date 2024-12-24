package store_fs

import (
	"encoding/gob"
	"strings"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/echo/query_spec"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
)

func init() {
	gob.Register(sku.Transacted{})
}

func Make(
	k sku.Config,
	dp interfaces.FuncIter[*fd.FD],
	fileExtensions interfaces.FileExtensionGetter,
	st dir_layout.DirLayout,
	ofo object_inventory_format.Options,
	fileEncoder FileEncoder,
) (fs *Store, err error) {
	fs = &Store{
		config:         k,
		deletedPrinter: dp,
		dirLayout:      st,
		fileEncoder:    fileEncoder,
		fileExtensions: fileExtensions,
		dir:            st.Cwd(),
		dirItems:       makeObjectsWithDir(st.Cwd(), fileExtensions, st),
		deleted: collections_value.MakeMutableValueSet[*fd.FD](
			nil,
		),
		deletedInternal: collections_value.MakeMutableValueSet[*fd.FD](
			nil,
		),
		objectFormatOptions: ofo,
		metadataTextParser: object_metadata.MakeTextParser(
			st,
			nil,
		),
	}

	return
}

type Store struct {
	config              sku.Config
	deletedPrinter      interfaces.FuncIter[*fd.FD]
	metadataTextParser  object_metadata.TextParser
	dirLayout           dir_layout.DirLayout
	fileEncoder         FileEncoder
	ic                  ids.InlineTypeChecker
	fileExtensions      interfaces.FileExtensionGetter
	dir                 string
	objectFormatOptions object_inventory_format.Options

	dirItems

	deleteLock      sync.Mutex
	deleted         fd.MutableSet
	deletedInternal fd.MutableSet
}

func (fs *Store) GetExternalStoreLike() external_store.StoreLike {
	return fs
}

// Deletions of user objects that should be exposed to the user
func (s *Store) DeleteCheckedOut(co *sku.CheckedOut) (err error) {
	external := co.GetSkuExternal()

	var i *sku.FSItem

	if i, err = s.ReadFSItemFromExternal(external); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.deleteLock.Lock()
	defer s.deleteLock.Unlock()

	if err = i.MutableSetLike.Each(s.deleted.Add); err != nil {
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
		fs.dirLayout,
		fs.deletedPrinter,
		fs.deleted,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = deleteOp.Run(
		fs.config.IsDryRun(),
		fs.dirLayout,
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
	sb.WriteRune(query_spec.OpGroupOpen)

	hasOne := false

	writeOneIfNecessary := func(v interfaces.Stringer) (err error) {
		if hasOne {
			sb.WriteRune(query_spec.OpOr)
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

	sb.WriteRune(query_spec.OpGroupClose)

	out = sb.String()
	return
}

func (s *Store) GetExternalObjectIds() (ks []sku.ExternalObjectId, err error) {
	if err = s.dirItems.processRootDir(); err != nil {
		err = errors.Wrap(err)
		return
	}

	ks = make([]sku.ExternalObjectId, 0)
	var l sync.Mutex

	if err = s.All(
		func(kfp *sku.FSItem) (err error) {
			l.Lock()
			defer l.Unlock()

			ks = append(ks, kfp.GetExternalObjectId())

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) GetObjectIdsForDir(
	fd *fd.FD,
) (k []sku.ExternalObjectId, err error) {
	if !fd.IsDir() {
		err = errors.Errorf("not a directory: %q", fd)
		return
	}

	var results []*sku.FSItem

	if results, err = s.dirItems.processDir(fd.GetPath()); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, r := range results {
		k = append(k, r.GetExternalObjectId())
	}

	return
}

// TODO confirm against actual Object Id
func (s *Store) GetObjectIdsForString(v string) (k []sku.ExternalObjectId, err error) {
	if v == "." {
		if k, err = s.GetExternalObjectIds(); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	var fdee *fd.FD

	if fdee, err = fd.MakeFromPath(v, s.dirLayout); err != nil {
		err = errors.Wrap(err)
		return
	}

	if fdee.IsDir() {
		if k, err = s.GetObjectIdsForDir(fdee); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		var results []*sku.FSItem

		if _, results, err = s.dirItems.processFD(fdee); err != nil {
			err = errors.Wrap(err)
			return
		}

		k = make([]sku.ExternalObjectId, 0, len(results))

		for _, r := range results {
			k = append(k, r.GetExternalObjectId())
		}
	}

	return
}

func (fs *Store) Get(
	k interfaces.ObjectId,
) (t *sku.FSItem, ok bool) {
	return fs.dirItems.probablyCheckedOut.Get(k.String())
}

func (s *Store) Initialize(esi external_store.Supplies) (err error) {
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

	if sk.ExternalObjectId.IsEmpty() {
		item.ExternalObjectId.ResetWith(&sk.ObjectId)
	} else {
		item.ExternalObjectId.ResetWith(&sk.ExternalObjectId)
	}

	return
}

func (s *Store) WriteFSItemToExternal(
	item *sku.FSItem,
	tg sku.TransactedGetter,
) (err error) {
	e := tg.GetSku()
	e.Metadata.Fields = e.Metadata.Fields[:0]

	m := &e.Metadata
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
		after := s.dirLayout.Rel(before)

		if err = e.ExternalObjectId.SetBlob(after); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		k := &item.ExternalObjectId

		e.ExternalObjectId.ResetWith(k)

		if e.ExternalObjectId.String() != k.String() {
			err = errors.Errorf("expected %q but got %q", k, &e.ExternalObjectId)
		}
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

		e.Metadata.Fields = append(e.Metadata.Fields, field)
	}

	return
}
