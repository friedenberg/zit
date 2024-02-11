package store_util

import (
	"code.linenisgreat.com/zit-go/src/alfa/errors"
	"code.linenisgreat.com/zit-go/src/charlie/collections"
	"code.linenisgreat.com/zit-go/src/charlie/sha"
	"code.linenisgreat.com/zit-go/src/delta/standort"
	"code.linenisgreat.com/zit-go/src/delta/thyme"
	"code.linenisgreat.com/zit-go/src/echo/kennung"
	"code.linenisgreat.com/zit-go/src/golf/ennui"
	"code.linenisgreat.com/zit-go/src/golf/kennung_index"
	"code.linenisgreat.com/zit-go/src/golf/objekte_format"
	"code.linenisgreat.com/zit-go/src/hotel/sku"
	"code.linenisgreat.com/zit-go/src/india/objekte_collections"
	"code.linenisgreat.com/zit-go/src/juliett/konfig"
	"code.linenisgreat.com/zit-go/src/kilo/cwd"
	"code.linenisgreat.com/zit-go/src/kilo/store_verzeichnisse"
	"code.linenisgreat.com/zit-go/src/lima/akten"
	"code.linenisgreat.com/zit-go/src/lima/bestandsaufnahme"
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

func (s *common) GetAkten() *akten.Akten {
	return s.akten
}

func (s *common) GetEnnui() ennui.Ennui {
	return nil
}

func (s *common) GetFileEncoder() objekte_collections.FileEncoder {
	return s.fileEncoder
}

func (s *common) GetCwdFiles() *cwd.CwdFiles {
	return s.cwdFiles
}

func (s *common) GetObjekteFormatOptions() objekte_format.Options {
	return s.options
}

func (s *common) GetPersistentMetadateiFormat() objekte_format.Format {
	return s.persistentMetadateiFormat
}

func (s *common) GetTime() thyme.Time {
	return thyme.Now()
}

func (s *common) GetTai() kennung.Tai {
	return kennung.NowTai()
}

func (s *common) GetBestandsaufnahmeStore() bestandsaufnahme.Store {
	return s.bestandsaufnahmeStore
}

func (s *common) GetAbbrStore() AbbrStore {
	return s.Abbr
}

func (s *common) GetKennungIndex() kennung_index.Index {
	return s.kennungIndex
}

func (s *common) GetStandort() standort.Standort {
	return s.standort
}

func (s *common) GetKonfig() *konfig.Compiled {
	return s.konfig
}

func (s *common) GetVerzeichnisse() *store_verzeichnisse.Store {
	return s.verzeichnisse
}

func (s *common) ReadOneEnnui(sh *sha.Sha) (*sku.Transacted, error) {
	if s.konfig.GetStoreVersion().GetInt() > 4 {
		return s.GetBestandsaufnahmeStore().ReadOneEnnui(sh)
	} else {
		return s.GetVerzeichnisse().ReadOneShas(sh)
	}
}

func (s *common) ReadOneKennung(k kennung.Kennung) (sk *sku.Transacted, err error) {
	if s.konfig.GetStoreVersion().GetInt() > 4 {
		return s.GetBestandsaufnahmeStore().ReadOneKennung(k)
	} else {
		return s.GetVerzeichnisse().ReadOneKennung(k)
	}
}

func (s *common) ReaderFor(sh *sha.Sha) (rc sha.ReadCloser, err error) {
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
