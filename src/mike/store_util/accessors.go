package store_util

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/delta/thyme"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/ennui"
	"github.com/friedenberg/zit/src/golf/kennung_index"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/objekte_collections"
	"github.com/friedenberg/zit/src/juliett/konfig"
	"github.com/friedenberg/zit/src/kilo/cwd"
	"github.com/friedenberg/zit/src/kilo/objekte_store"
	"github.com/friedenberg/zit/src/kilo/store_verzeichnisse"
	"github.com/friedenberg/zit/src/lima/akten"
	"github.com/friedenberg/zit/src/lima/bestandsaufnahme"
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

func (s *common) ReaderFor(sh *sha.Sha) (rc sha.ReadCloser, err error) {
	if rc, err = s.standort.AkteReaderFrom(
		sh,
		s.standort.DirVerzeichnisseMetadatei(),
	); err != nil {
		if errors.IsNotExist(err) {
			err = objekte_store.ErrNotFoundEmpty
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	return
}
