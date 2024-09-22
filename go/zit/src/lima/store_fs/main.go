package store_fs

import (
	"encoding/gob"
	"strings"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_extensions"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
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

type Item = sku.FSItem

type Store struct {
	config              sku.Config
	deletedPrinter      interfaces.FuncIter[*fd.FD]
	metadataTextParser  object_metadata.TextParser
	fs_home             fs_home.Home
	fileEncoder         FileEncoder
	ic                  ids.InlineTypeChecker
	fileExtensions      file_extensions.FileExtensions
	dir                 string
	objectFormatOptions object_inventory_format.Options

	dirItems

	deleteLock sync.Mutex
	deleted    fd.MutableSet
}

func (fs *Store) GetExternalStoreLike() external_store.StoreLike {
	return fs
}

func (s *Store) DeleteExternalLike(el sku.ExternalLike) (err error) {
	e := el.(*sku.Transacted)

	var i *Item

	if i, err = s.ReadFSItemFromExternal(e); err != nil {
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

func (fs *Store) Flush() (err error) {
	deleteOp := DeleteCheckout{}

	if err = deleteOp.Run(
		fs.config.IsDryRun(),
		fs.fs_home,
		fs.deletedPrinter,
		fs.deleted,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	fs.deleted.Reset()

	return
}

func (fs *Store) String() (out string) {
	if iter.Len(
		fs.dirItems.objects,
		fs.blobs,
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

	fs.dirItems.objects.Each(
		func(z *Item) (err error) {
			return writeOneIfNecessary(z)
		},
	)

	fs.blobs.Each(
		func(z *Item) (err error) {
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
		func(kfp *Item) (err error) {
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

func (s *Store) GetObjectIdsForDir(
	fd *fd.FD,
) (k []sku.ExternalObjectId, err error) {
	if !fd.IsDir() {
		err = errors.Errorf("not a directory: %q", fd)
		return
	}

	var results []*Item

	if results, err = s.dirItems.processDir(fd.GetPath()); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, r := range results {
		k = append(k, r)
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

	if fdee, err = fd.MakeFromPath(v, s.fs_home); err != nil {
		err = errors.Wrap(err)
		return
	}

	if fdee.IsDir() {
		if k, err = s.GetObjectIdsForDir(fdee); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		var results []*Item

		if _, results, err = s.dirItems.processFD(fdee); err != nil {
			err = errors.Wrap(err)
			return
		}

		k = make([]sku.ExternalObjectId, 0, len(results))

		for _, r := range results {
			k = append(k, r)
		}
	}

	return
}

func (fs *Store) Get(
	k interfaces.ObjectId,
) (t *Item, ok bool) {
	return fs.dirItems.objects.Get(k.String())
}

func (s *Store) Initialize(esi external_store.Supplies) (err error) {
	s.externalStoreSupplies = esi
	return
}

func (s *Store) ApplyDotOperator() (err error) {
	if err = s.dirItems.processRootDir(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadFSItemFromExternal(el sku.ExternalLike) (i *Item, err error) {
	i = &Item{} // TODO use pool or use dir_items?
	i.Reset()

	e := el.(*sku.Transacted)

	// TODO handle sort order
	for _, f := range e.Metadata.Fields {
		var fdee *fd.FD
		switch strings.ToLower(f.Key) {
		case "object":
			fdee = &i.Object

		case "blob":
			fdee = &i.Blob

		case "conflict":
			fdee = &i.Conflict

		default:
			err = errors.Errorf("unexpected field: %#v", f)
			return
		}

		// if we've already set one of object, blob, or conflict, don't set it again
		// and instead add a new FD to the item
		if !fdee.IsEmpty() {
			fdee = &fd.FD{}
		}

		if err = fdee.SetIgnoreNotExists(f.Value); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = i.Add(fdee); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	i.State = e.State
	i.ObjectId.ResetWith(&e.ExternalObjectId)

	return
}

func (s *Store) WriteFSItemToExternal(i *Item, el sku.ExternalLike) (err error) {
	e := el.(*sku.Transacted)
	e.Metadata.Fields = e.Metadata.Fields[:0]
	k := &i.ObjectId

	e.ExternalObjectId.ResetWith(k)

	if e.ExternalObjectId.String() != k.String() {
		err = errors.Errorf("expected %q but got %q", k, &e.ExternalObjectId)
	}

	m := &e.Metadata
	m.Tai = i.GetTai()

	fdees := iter.SortedValues(i.MutableSetLike)

	for _, f := range fdees {
		field := object_metadata.Field{
			Value:     f.GetPath(),
			ColorType: string_format_writer.ColorTypeId,
		}

		switch {
		case i.Object.Equals(f):
			field.Key = "object"

		case i.Conflict.Equals(f):
			field.Key = "conflict"

		case i.Blob.Equals(f):
			fallthrough

		default:
			field.Key = "blob"
		}

		e.Metadata.Fields = append(e.Metadata.Fields, field)
	}

	return
}
