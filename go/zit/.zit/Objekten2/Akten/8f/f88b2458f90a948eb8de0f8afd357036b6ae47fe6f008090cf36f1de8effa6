package store

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/thyme"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/zettel_id_index"
	"code.linenisgreat.com/zit/go/zit/src/golf/mutable_config_blobs"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_probe_index"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/blob_store"
	"code.linenisgreat.com/zit/go/zit/src/juliett/config"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
	"code.linenisgreat.com/zit/go/zit/src/kilo/stream_index"
	"code.linenisgreat.com/zit/go/zit/src/lima/inventory_list_store"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_fs"
)

func (u *Store) GetBrowserStore() *external_store.Store {
	return u.externalStores[*(ids.MustRepoId("browser"))]
}

func (s *Store) GetBlobStore() *blob_store.VersionedStores {
	return s.blobStore
}

func (s *Store) GetEnnui() object_probe_index.Index {
	return nil
}

func (s *Store) GetStoreFS() *store_fs.Store {
	return s.storeFS
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

func (s *Store) GetAbbrStore() AbbrStore {
	return s.Abbr
}

func (s *Store) GetZettelIdIndex() zettel_id_index.Index {
	return s.zettelIdIndex
}

func (s *Store) GetDirectoryLayout() dir_layout.DirLayout {
	return s.dirLayout
}

func (s *Store) GetConfig() *config.Compiled {
	return s.config
}

func (s *Store) GetStreamIndex() *stream_index.Index {
	return s.streamIndex
}

func (s *Store) GetConfigBlobFormat() blob_store.Format2[mutable_config_blobs.Blob] {
	return s.configBlobFormat
}
