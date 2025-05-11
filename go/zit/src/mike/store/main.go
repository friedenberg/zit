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
	"code.linenisgreat.com/zit/go/zit/src/lima/typed_blob_store"
	"code.linenisgreat.com/zit/go/zit/src/mike/store_config"
	"code.linenisgreat.com/zit/go/zit/src/mike/store_workspace"
)

type Store struct {
	sunrise      ids.Tai
	config       store_config.StoreMutable
	envRepo      env_repo.Env
	envWorkspace env_workspace.Env

	typedBlobStore     typed_blob_store.Stores
	inventoryListStore inventory_list_store.Store
	Abbr               sku.AbbrStore

	inventoryList          *sku.OpenList
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

func (store *Store) Initialize(
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
	store.config = config
	store.envRepo = envRepo
	store.envWorkspace = envWorkspace
	store.typedBlobStore = typedBlobStore
	store.persistentObjectFormat = pmf
	store.sunrise = sunrise
	store.envLua = envLua
	store.queryBuilder = queryBuilder
	store.dormantIndex = dormantIndex

	store.Abbr = abbrStore

	if err = store.inventoryListStore.Initialize(
		store.GetEnvRepo(),
		store,
		typedBlobStore.InventoryList,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if store.inventoryList, err = store.inventoryListStore.MakeOpenList(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if store.zettelIdIndex, err = zettel_id_index.MakeIndex(
		// TODO
		store.GetConfig(),
		store.GetEnvRepo(),
		store.GetEnvRepo(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if store.streamIndex, err = stream_index.MakeIndex(
		store.GetEnvRepo(),
		store.applyDormantAndRealizeTags,
		store.GetEnvRepo().DirCacheObjects(),
		store.sunrise,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	store.protoZettel = sku.MakeProto(
		store.envWorkspace.GetDefaults(),
	)

	store.configBlobFormat = typed_blob_store.MakeBlobFormat2(
		typed_blob_store.MakeTextParserIgnoreTomlErrors2[config_mutable_blobs.Blob](
			store.GetEnvRepo(),
		),
		typed_blob_store.ParsedBlobTomlFormatter2[config_mutable_blobs.Blob]{},
		store.GetEnvRepo(),
	)

	return
}

func (store *Store) MakeSupplies(
	repoId ids.RepoId,
) (supplies store_workspace.Supplies) {
	supplies.WorkspaceDir = store.envWorkspace.GetWorkspaceDir()
	supplies.ObjectStore = store

	supplies.Env = store.GetEnvRepo()
	supplies.Clock = store.sunrise
	supplies.BlobStore = store.typedBlobStore
	supplies.RepoId = repoId
	supplies.DirCache = store.GetEnvRepo().DirCacheRepo(repoId.GetRepoIdString())

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
