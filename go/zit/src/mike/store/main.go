package store

import (
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/zettel_id_index"
	"code.linenisgreat.com/zit/go/zit/src/golf/config_mutable_blobs"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
	"code.linenisgreat.com/zit/go/zit/src/hotel/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/box_format"
	"code.linenisgreat.com/zit/go/zit/src/kilo/dormant_index"
	"code.linenisgreat.com/zit/go/zit/src/kilo/env_workspace"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/stream_index"
	"code.linenisgreat.com/zit/go/zit/src/lima/env_lua"
	"code.linenisgreat.com/zit/go/zit/src/lima/inventory_list_store"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/lima/typed_blob_store"
	"code.linenisgreat.com/zit/go/zit/src/mike/store_config"
	"code.linenisgreat.com/zit/go/zit/src/mike/store_workspace"
)

type Store struct {
	sunrise      ids.Tai
	config       store_config.StoreMutable
	envRepo      env_repo.Env
	envWorkspace env_workspace.Env

	externalStores     map[ids.RepoId]*env_workspace.Store
	typedBlobStore     typed_blob_store.Stores
	inventoryListStore inventory_list_store.Store
	Abbr               sku.AbbrStore

	inventoryList          *sku.List
	persistentObjectFormat object_inventory_format.Format
	configBlobFormat       interfaces.Format[config_mutable_blobs.Blob]
	envLua                 env_lua.Env
	tagLock                sync.Mutex

	streamIndex   *stream_index.Index
	zettelIdIndex zettel_id_index.Index
	dormantIndex  *dormant_index.Index

	protoZettel  sku.Proto
	queryBuilder *query.Builder

	ui sku.UIStorePrinters
}

func (c *Store) Initialize(
	config store_config.StoreMutable,
	envRepo env_repo.Env,
	envWorkspace env_workspace.Env,
	pmf object_inventory_format.Format,
	sunrise ids.Tai,
	envLua env_lua.Env,
	queryBuilder *query.Builder,
	box *box_format.BoxTransacted,
	typedBlobStore typed_blob_store.Stores,
	dormantIndex *dormant_index.Index,
	abbrStore sku.AbbrStore,
) (err error) {
	c.config = config
	c.envRepo = envRepo
	c.envWorkspace = envWorkspace
	c.typedBlobStore = typedBlobStore
	c.persistentObjectFormat = pmf
	c.sunrise = sunrise
	c.envLua = envLua
	c.queryBuilder = queryBuilder
	c.dormantIndex = dormantIndex

	c.inventoryList = sku.MakeList()

	c.Abbr = abbrStore

	if err = c.inventoryListStore.Initialize(
		c.GetEnvRepo(),
		c,
		typedBlobStore.InventoryList,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if c.zettelIdIndex, err = zettel_id_index.MakeIndex(
		// TODO
		c.GetConfig(),
		c.GetEnvRepo(),
		c.GetEnvRepo(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if c.streamIndex, err = stream_index.MakeIndex(
		c.GetEnvRepo(),
		c.applyDormantAndRealizeTags,
		c.GetEnvRepo().DirCacheObjects(),
		c.sunrise,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	c.protoZettel = sku.MakeProto(
		c.envWorkspace.GetDefaults(),
	)

	c.configBlobFormat = typed_blob_store.MakeBlobFormat2(
		typed_blob_store.MakeTextParserIgnoreTomlErrors2[config_mutable_blobs.Blob](
			c.GetEnvRepo(),
		),
		typed_blob_store.ParsedBlobTomlFormatter2[config_mutable_blobs.Blob]{},
		c.GetEnvRepo(),
	)

	return
}

// TODO add external_store.Supplies to Store and just use that
func (store *Store) MakeSupplies() (supplies store_workspace.Supplies) {
	supplies.WorkspaceDir = store.envWorkspace.GetWorkspaceDir()
	supplies.ObjectStore = store

	supplies.Env = store.GetEnvRepo()
	supplies.Clock = store.sunrise
	supplies.BlobStore = store.typedBlobStore

	return
}

func (s *Store) SetExternalStores(
	stores map[ids.RepoId]*env_workspace.Store,
) (err error) {
	s.externalStores = stores

	supplies := s.MakeSupplies()

	for k, es := range s.externalStores {
		supplies.RepoId = k
		supplies.DirCache = s.GetEnvRepo().DirCacheRepo(k.GetRepoIdString())
		es.Supplies = supplies

		if _, ok := es.StoreLike.(*store_fs.Store); ok {
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

func (s *Store) SetUIDelegate(ud sku.UIStorePrinters) {
	s.ui = ud
	s.inventoryListStore.SetUIDelegate(ud)
}
