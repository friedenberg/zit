package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/echo/standort"
	"code.linenisgreat.com/zit/go/zit/src/echo/thyme"
	"code.linenisgreat.com/zit/go/zit/src/external_store"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/erworben"
	"code.linenisgreat.com/zit/go/zit/src/golf/ennui"
	"code.linenisgreat.com/zit/go/zit/src/golf/kennung_index"
	"code.linenisgreat.com/zit/go/zit/src/golf/objekte_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/akten"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/juliett/konfig"
	"code.linenisgreat.com/zit/go/zit/src/kilo/store_verzeichnisse"
	"code.linenisgreat.com/zit/go/zit/src/kilo/zettel"
	"code.linenisgreat.com/zit/go/zit/src/lima/bestandsaufnahme"
)

func (u *Store) GetChrestStore() *external_store.Store {
	return u.externalStores["chrome"]
}

func (s *Store) GetAkten() *akten.Akten {
	return s.akten
}

func (s *Store) GetEnnui() ennui.Ennui {
	return nil
}

func (s *Store) GetFileEncoder() store_fs.FileEncoder {
	return s.fileEncoder
}

func (s *Store) GetCwdFiles() *store_fs.Store {
	return s.cwdFiles
}

func (s *Store) GetObjekteFormatOptions() objekte_format.Options {
	return s.options
}

func (s *Store) GetProtoZettel() zettel.ProtoZettel {
	return s.protoZettel
}

func (s *Store) GetPersistentMetadateiFormat() objekte_format.Format {
	return s.persistentMetadateiFormat
}

func (s *Store) GetTime() thyme.Time {
	return thyme.Now()
}

func (s *Store) GetTai() kennung.Tai {
	return kennung.NowTai()
}

func (s *Store) GetBestandsaufnahmeStore() bestandsaufnahme.Store {
	return s.bestandsaufnahmeStore
}

func (s *Store) GetAbbrStore() AbbrStore {
	return s.Abbr
}

func (s *Store) GetKennungIndex() kennung_index.Index {
	return s.kennungIndex
}

func (s *Store) GetStandort() standort.Standort {
	return s.standort
}

func (s *Store) GetKonfig() *konfig.Compiled {
	return s.konfig
}

func (s *Store) GetVerzeichnisse() *store_verzeichnisse.Store {
	return s.verzeichnisse
}

func (s *Store) GetKonfigAkteFormat() akten.Format[erworben.Akte, *erworben.Akte] {
	return s.konfigAkteFormat
}

func (s *Store) ReadOneEnnui(sh *sha.Sha) (*sku.Transacted, error) {
	return s.GetVerzeichnisse().ReadOneEnnui(sh)
}

func (s *Store) ReadOneKennung(
	k schnittstellen.StringerGattungGetter,
) (sk *sku.Transacted, err error) {
	return s.GetVerzeichnisse().ReadOneKennung(k)
}

func (s *Store) ReaderFor(sh *sha.Sha) (rc sha.ReadCloser, err error) {
	if rc, err = s.standort.AkteReaderFrom(
		sh,
		s.standort.DirVerzeichnisseMetadateiKennungMutter(),
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
