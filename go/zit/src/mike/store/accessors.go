package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/thyme"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/zettel_id_index"
	"code.linenisgreat.com/zit/go/zit/src/golf/config_mutable_blobs"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
	"code.linenisgreat.com/zit/go/zit/src/hotel/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/india/object_probe_index"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/env_workspace"
	"code.linenisgreat.com/zit/go/zit/src/kilo/stream_index"
	"code.linenisgreat.com/zit/go/zit/src/lima/inventory_list_store"
	"code.linenisgreat.com/zit/go/zit/src/lima/typed_blob_store"
	"code.linenisgreat.com/zit/go/zit/src/mike/store_config"
)

func (u *Store) GetBrowserStore() *env_workspace.Store {
	return u.externalStores[*(ids.MustRepoId("browser"))]
}

func (s *Store) GetTypedBlobStore() typed_blob_store.Stores {
	return s.typedBlobStore
}

func (s *Store) GetEnnui() object_probe_index.Index {
	return nil
}

func (s *Store) GetProtoZettel() sku.Proto {
	return s.protoZettel
}

func (s *Store) GetPersistentMetadataFormat() object_inventory_format.Format {
	return s.persistentObjectFormat
}

func (s *Store) GetTime() thyme.Time {
	return thyme.Now()
}

func (s *Store) GetTai() ids.Tai {
	return ids.NowTai()
}

func (s *Store) GetInventoryListStore() *inventory_list_store.Store {
	return &s.inventoryListStore
}

func (s *Store) GetAbbrStore() sku.AbbrStore {
	return s.Abbr
}

func (s *Store) GetZettelIdIndex() zettel_id_index.Index {
	return s.zettelIdIndex
}

func (s *Store) GetEnvRepo() env_repo.Env {
	return s.envRepo
}

func (s *Store) GetConfig() store_config.Store {
	return s.config
}

func (s *Store) GetStreamIndex() *stream_index.Index {
	return s.streamIndex
}

func (s *Store) GetConfigBlobFormat() interfaces.Format[config_mutable_blobs.Blob] {
	return s.configBlobFormat
}
