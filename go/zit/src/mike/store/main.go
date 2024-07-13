package store

import (
	"flag"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/echo/standort"
	"code.linenisgreat.com/zit/go/zit/src/echo/thyme"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/erworben"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/go/zit/src/golf/kennung_index"
	"code.linenisgreat.com/zit/go/zit/src/golf/objekte_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/akten"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/juliett/konfig"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
	"code.linenisgreat.com/zit/go/zit/src/kilo/store_verzeichnisse"
	"code.linenisgreat.com/zit/go/zit/src/kilo/zettel"
	"code.linenisgreat.com/zit/go/zit/src/lima/bestandsaufnahme"
)

type Store struct {
	konfig                    *konfig.Compiled
	standort                  standort.Standort
	cwdFiles                  *store_fs.Store
	externalStores            map[string]*external_store.Store
	akten                     *akten.Akten
	bestandsaufnahmeAkte      bestandsaufnahme.InventoryList
	options                   objekte_format.Options
	Abbr                      AbbrStore
	persistentMetadateiFormat objekte_format.Format
	fileEncoder               store_fs.FileEncoder
	luaVMPoolBuilder          *lua.VMPoolBuilder
	etikettenLock             sync.Mutex

	verzeichnisse *store_verzeichnisse.Store

	sonnenaufgang thyme.Time

	checkedOutLogPrinter interfaces.FuncIter[sku.CheckedOutLike]

	metadateiTextParser metadatei.TextParser

	bestandsaufnahmeStore bestandsaufnahme.Store
	kennungIndex          kennung_index.Index

	sku.TransactedAdder
	typenIndex kennung_index.KennungIndex[kennung.Typ, *kennung.Typ]

	protoZettel      zettel.ProtoZettel
	konfigAkteFormat akten.Format[erworben.Akte, *erworben.Akte]

	queryBuilder *query.Builder

	Logger
}

type Logger struct {
	New, Updated, Unchanged interfaces.FuncIter[*sku.Transacted]
}

func (c *Store) Initialize(
	flags *flag.FlagSet,
	k *konfig.Compiled,
	st standort.Standort,
	pmf objekte_format.Format,
	t thyme.Time,
	luaVMPoolBuilder *lua.VMPoolBuilder,
	qb *query.Builder,
	options objekte_format.Options,
) (err error) {
	c.konfig = k
	c.standort = st
	c.akten = akten.Make(st)
	c.persistentMetadateiFormat = pmf
	c.options = options
	c.sonnenaufgang = t
	c.fileEncoder = store_fs.MakeFileEncoder(st, k)
	c.luaVMPoolBuilder = luaVMPoolBuilder
	c.queryBuilder = qb

	c.metadateiTextParser = metadatei.MakeTextParser(
		c.standort,
		nil, // TODO-P1 make akteFormatter
	)

	c.typenIndex = kennung_index.MakeIndex2[kennung.Typ](
		c.standort,
		st.DirVerzeichnisse("TypenIndexV0"),
	)

	c.bestandsaufnahmeAkte = bestandsaufnahme.InventoryList{
		Skus: sku.MakeTransactedHeap(),
	}

	if c.Abbr, err = newIndexAbbr(
		c.standort,
		st.DirVerzeichnisse("Abbr"),
	); err != nil {
		err = errors.Wrapf(err, "failed to init abbr index")
		return
	}

	if c.bestandsaufnahmeStore, err = bestandsaufnahme.MakeStore(
		c.GetStandort(),
		c.GetStandort().GetLockSmith(),
		c.konfig.GetStoreVersion(),
		c.standort.ObjekteReaderWriterFactory(gattung.Bestandsaufnahme),
		c.standort,
		pmf,
		c,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if c.kennungIndex, err = kennung_index.MakeIndex(
		c.GetKonfig(),
		c.GetStandort(),
		c.GetStandort(),
	); err != nil {
		err = errors.Wrapf(err, "failed to init zettel index")
		return
	}

	if c.verzeichnisse, err = store_verzeichnisse.MakeStore(
		c.GetStandort(),
		c.GetKonfig(),
		c.GetStandort().DirVerzeichnisseObjekten(),
		pmf,
		c.options,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	c.protoZettel = zettel.MakeProtoZettel(c.GetKonfig())

	c.konfigAkteFormat = akten.MakeAkteFormat(
		akten.MakeTextParserIgnoreTomlErrors[erworben.Akte](
			c.GetStandort(),
		),
		akten.ParsedAkteTomlFormatter[erworben.Akte, *erworben.Akte]{},
		c.GetStandort(),
	)

	return
}

func (s *Store) SetExternalStores(
	stores map[string]*external_store.Store,
) (err error) {
	s.externalStores = stores

	for k, es := range s.externalStores {
		es.StoreFuncs = external_store.StoreFuncs{
			FuncRealize:     s.tryRealize,
			FuncCommit:      s.tryRealizeAndOrStore,
			FuncReadSha:     s.ReadOneEnnui,
			FuncReadOneInto: s.ReadOneInto,
			FuncQuery:       s.Query,
		}

		es.Standort = s.GetStandort()
		es.DirCache = s.GetStandort().DirVerzeichnisseKasten(k)

		if esfs, ok := es.StoreLike.(*store_fs.Store); ok {
			s.cwdFiles = esfs

			// TODO remove once store_fs.Store is fully ExternalStoreLike
			if err = es.Initialize(); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	return
}

// TODO remove
func (s *Store) SetCheckedOutLogWriter(
	zelw interfaces.FuncIter[sku.CheckedOutLike],
) {
	s.checkedOutLogPrinter = zelw
}

func (s *Store) ResetIndexes() (err error) {
	if err = s.typenIndex.Reset(); err != nil {
		err = errors.Wrapf(err, "failed to reset etiketten index")
		return
	}

	if err = s.kennungIndex.Reset(); err != nil {
		err = errors.Wrapf(err, "failed to reset index kennung")
		return
	}

	return
}

func (s *Store) SetLogWriter(lw Logger) {
	s.Logger = lw
}
