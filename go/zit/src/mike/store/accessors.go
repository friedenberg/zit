package store

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/src/charlie/sha"
	"code.linenisgreat.com/zit/src/delta/standort"
	"code.linenisgreat.com/zit/src/delta/thyme"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/golf/ennui"
	"code.linenisgreat.com/zit/src/golf/kennung_index"
	"code.linenisgreat.com/zit/src/golf/objekte_format"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/erworben"
	"code.linenisgreat.com/zit/src/india/objekte_collections"
	"code.linenisgreat.com/zit/src/juliett/konfig"
	"code.linenisgreat.com/zit/src/juliett/objekte"
	"code.linenisgreat.com/zit/src/kilo/cwd"
	"code.linenisgreat.com/zit/src/kilo/store_verzeichnisse"
	"code.linenisgreat.com/zit/src/lima/akten"
	"code.linenisgreat.com/zit/src/lima/bestandsaufnahme"
)

type accessors interface {
	standort.Getter
	konfig.Getter
	objekte_format.Getter

	GetAbbrStore() AbbrStore
	GetAkten() *akten.Akten
	GetBestandsaufnahmeStore() bestandsaufnahme.Store
	GetCwdFiles() *cwd.CwdFiles
	GetEnnui() ennui.Ennui
	GetFileEncoder() objekte_collections.FileEncoder
	GetKennungIndex() kennung_index.Index
	GetObjekteFormatOptions() objekte_format.Options
	GetVerzeichnisse() *store_verzeichnisse.Store
	ReadOneEnnui(*sha.Sha) (*sku.Transacted, error)
	ReadOneKennung(kennung.Kennung) (*sku.Transacted, error)
	ReaderFor(*sha.Sha) (sha.ReadCloser, error)
}

func (s *Store) GetAkten() *akten.Akten {
	return s.akten
}

func (s *Store) GetEnnui() ennui.Ennui {
	return nil
}

func (s *Store) GetFileEncoder() objekte_collections.FileEncoder {
	return s.fileEncoder
}

func (s *Store) GetCwdFiles() *cwd.CwdFiles {
	return s.cwdFiles
}

func (s *Store) GetObjekteFormatOptions() objekte_format.Options {
	return s.options
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

func (s *Store) GetKonfigAkteFormat() objekte.AkteFormat[erworben.Akte, *erworben.Akte] {
	return s.konfigAkteFormat
}

func (s *Store) ReadOneEnnui(sh *sha.Sha) (*sku.Transacted, error) {
	if s.konfig.GetStoreVersion().GetInt() > 4 {
		return s.GetBestandsaufnahmeStore().ReadOneEnnui(sh)
	} else {
		return s.GetVerzeichnisse().ReadOneShas(sh)
	}
}

func (s *Store) ReadOneKennung(k kennung.Kennung) (sk *sku.Transacted, err error) {
	if s.konfig.GetStoreVersion().GetInt() > 4 {
		return s.GetBestandsaufnahmeStore().ReadOneKennung(k)
	} else {
		return s.GetVerzeichnisse().ReadOneKennung(k)
	}
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
