package store

import (
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/zettel_id_index"
	"code.linenisgreat.com/zit/go/zit/src/golf/config_mutable_blobs"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
	"code.linenisgreat.com/zit/go/zit/src/hotel/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/box_format"
	"code.linenisgreat.com/zit/go/zit/src/kilo/dormant_index"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/stream_index"
	"code.linenisgreat.com/zit/go/zit/src/lima/blob_store"
	"code.linenisgreat.com/zit/go/zit/src/lima/inventory_list_store"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/mike/store_config"
)

type Store struct {
	sunrise   ids.Tai
	config    store_config.StoreMutable
	dirLayout env_repo.Env

	storeFS            *store_fs.Store
	externalStores     map[ids.RepoId]*external_store.Store
	blobStore          *blob_store.VersionedStores
	inventoryListStore inventory_list_store.Store
	Abbr               AbbrStore

	inventoryList          *sku.List
	options                object_inventory_format.Options
	persistentObjectFormat object_inventory_format.Format
	configBlobFormat       interfaces.Format[config_mutable_blobs.Blob]
	luaVMPoolBuilder       *lua.VMPoolBuilder
	tagLock                sync.Mutex

	streamIndex   *stream_index.Index
	zettelIdIndex zettel_id_index.Index
	dormantIndex  *dormant_index.Index

	protoZettel  sku.Proto
	queryBuilder *query.Builder

	ui UIDelegate
}

type UIDelegate struct {
	TransactedNew       interfaces.FuncIter[*sku.Transacted]
	TransactedUpdated   interfaces.FuncIter[*sku.Transacted]
	TransactedUnchanged interfaces.FuncIter[*sku.Transacted]

	CheckedOutCheckedOut interfaces.FuncIter[sku.SkuType]
	CheckedOutChanged    interfaces.FuncIter[sku.SkuType]
}

func (c *Store) Initialize(
	config store_config.StoreMutable,
	envRepo env_repo.Env,
	pmf object_inventory_format.Format,
	sunrise ids.Tai,
	luaVMPoolBuilder *lua.VMPoolBuilder,
	queryBuilder *query.Builder,
	options object_inventory_format.Options,
	box *box_format.BoxTransacted,
	blobStore *blob_store.VersionedStores,
	dormantIndex *dormant_index.Index,
) (err error) {
	c.config = config
	c.dirLayout = envRepo
	c.blobStore = blobStore
	c.persistentObjectFormat = pmf
	c.options = options
	c.sunrise = sunrise
	c.luaVMPoolBuilder = luaVMPoolBuilder
	c.queryBuilder = queryBuilder
	c.dormantIndex = dormantIndex

	c.inventoryList = sku.MakeList()

	if c.Abbr, err = newIndexAbbr(
		config.GetCLIConfig().PrintOptions,
		c.dirLayout,
		envRepo.DirCache("Abbr"),
	); err != nil {
		err = errors.Wrapf(err, "failed to init abbr index")
		return
	}

	if err = c.inventoryListStore.Initialize(
		c.GetDirectoryLayout(),
		pmf,
		c,
		blobStore.GetInventoryList(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if c.zettelIdIndex, err = zettel_id_index.MakeIndex(
		// TODO
		c.GetConfig(),
		c.GetDirectoryLayout(),
		c.GetDirectoryLayout(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if c.streamIndex, err = stream_index.MakeIndex(
		c.GetDirectoryLayout(),
		c.applyDormantAndRealizeTags,
		c.GetDirectoryLayout().DirCacheObjects(),
		c.sunrise,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	c.protoZettel = sku.MakeProto(
		config.GetMutableConfig().GetDefaults().GetType(),
		config.GetMutableConfig().GetDefaults().GetTags(),
	)

	c.configBlobFormat = blob_store.MakeBlobFormat2(
		blob_store.MakeTextParserIgnoreTomlErrors2[config_mutable_blobs.Blob](
			c.GetDirectoryLayout(),
		),
		blob_store.ParsedBlobTomlFormatter2[config_mutable_blobs.Blob]{},
		c.GetDirectoryLayout(),
	)

	return
}

func (s *Store) SetExternalStores(
	stores map[ids.RepoId]*external_store.Store,
) (err error) {
	s.externalStores = stores

	for k, es := range s.externalStores {
		es.StoreFuncs = external_store.StoreFuncs{
			FuncRealize:        s.tryRealize,
			FuncCommit:         s.tryRealizeAndOrStore,
			FuncReadOneInto:    s.GetStreamIndex().ReadOneObjectId,
			FuncPrimitiveQuery: s.GetStreamIndex().ReadPrimitiveQuery,
		}

		es.Env = s.GetDirectoryLayout()
		es.DirCache = s.GetDirectoryLayout().DirCacheRepo(k.GetRepoIdString())

		es.RepoId = k
		es.Clock = s.sunrise
		es.BlobStore = s.blobStore

		if esfs, ok := es.StoreLike.(*store_fs.Store); ok {
			s.storeFS = esfs

			// TODO remove once store_fs.Store is fully ExternalStoreLike
			if err = es.Initialize(); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	return
}

func (s *Store) ResetIndexes() (err error) {
	if err = s.zettelIdIndex.Reset(); err != nil {
		err = errors.Wrapf(err, "failed to reset index object id index")
		return
	}

	return
}

func (s *Store) SetUIDelegate(ud UIDelegate) {
	s.ui = ud
}
