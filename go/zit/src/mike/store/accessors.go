package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/delta/thyme"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/mutable_config"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_id_index"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_probe_index"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/blob_store"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/juliett/konfig"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
	"code.linenisgreat.com/zit/go/zit/src/kilo/store_verzeichnisse"
	"code.linenisgreat.com/zit/go/zit/src/kilo/zettel"
	"code.linenisgreat.com/zit/go/zit/src/lima/bestandsaufnahme"
)

func (u *Store) GetChrestStore() *external_store.Store {
	return u.externalStores["chrome"]
}

func (s *Store) GetAkten() *blob_store.VersionedStores {
	return s.blob_store
}

func (s *Store) GetEnnui() object_probe_index.Ennui {
	return nil
}

func (s *Store) GetFileEncoder() store_fs.FileEncoder {
	return s.fileEncoder
}

func (s *Store) GetCwdFiles() *store_fs.Store {
	return s.cwdFiles
}

func (s *Store) GetObjekteFormatOptions() object_inventory_format.Options {
	return s.options
}

func (s *Store) GetProtoZettel() zettel.ProtoZettel {
	return s.protoZettel
}

func (s *Store) GetPersistentMetadateiFormat() object_inventory_format.Format {
	return s.persistentMetadateiFormat
}

func (s *Store) GetTime() thyme.Time {
	return thyme.Now()
}

func (s *Store) GetTai() ids.Tai {
	return ids.NowTai()
}

func (s *Store) GetBestandsaufnahmeStore() bestandsaufnahme.Store {
	return s.bestandsaufnahmeStore
}

func (s *Store) GetAbbrStore() AbbrStore {
	return s.Abbr
}

func (s *Store) GetKennungIndex() object_id_index.Index {
	return s.kennungIndex
}

func (s *Store) GetStandort() fs_home.Home {
	return s.fs_home
}

func (s *Store) GetKonfig() *konfig.Compiled {
	return s.konfig
}

func (s *Store) GetVerzeichnisse() *store_verzeichnisse.Store {
	return s.verzeichnisse
}

func (s *Store) GetKonfigAkteFormat() blob_store.Format[mutable_config.Blob, *mutable_config.Blob] {
	return s.konfigAkteFormat
}

func (s *Store) ReadOneEnnui(sh *sha.Sha) (*sku.Transacted, error) {
	return s.GetVerzeichnisse().ReadOneObjectSha(sh)
}

func (s *Store) ReadOneKennung(
	k interfaces.ObjectId,
) (sk *sku.Transacted, err error) {
	return s.GetVerzeichnisse().ReadOneObjectId(k)
}

func (s *Store) ReaderFor(sh *sha.Sha) (rc sha.ReadCloser, err error) {
	if rc, err = s.fs_home.BlobReaderFrom(
		sh,
		s.fs_home.DirVerzeichnisseMetadateiKennungMutter(),
	); err != nil {
		if errors.IsNotExist(err) {
			err = collections.MakeErrNotFound(sh)
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	return
}
