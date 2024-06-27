package store

import (
	"flag"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
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
	"code.linenisgreat.com/zit/go/zit/src/kilo/store_verzeichnisse"
	"code.linenisgreat.com/zit/go/zit/src/kilo/zettel"
	"code.linenisgreat.com/zit/go/zit/src/lima/bestandsaufnahme"
)

type Store struct {
	konfig                    *konfig.Compiled
	standort                  standort.Standort
	cwdFiles                  *store_fs.Store
	externalStores            map[string]sku.ExternalStore
	akten                     *akten.Akten
	bestandsaufnahmeAkte      bestandsaufnahme.Akte
	options                   objekte_format.Options
	Abbr                      AbbrStore
	persistentMetadateiFormat objekte_format.Format
	fileEncoder               store_fs.FileEncoder
	virtualStores             map[string]*query.VirtualStoreInitable
	luaVMPoolBuilder          *lua.VMPoolBuilder
	etikettenLock             sync.Mutex

	verzeichnisse *store_verzeichnisse.Store

	sonnenaufgang thyme.Time

	checkedOutLogPrinter schnittstellen.FuncIter[sku.CheckedOutLike]

	metadateiTextParser metadatei.TextParser

	bestandsaufnahmeStore bestandsaufnahme.Store
	kennungIndex          kennung_index.Index

	sku.TransactedAdder
	typenIndex kennung_index.KennungIndex[kennung.Typ, *kennung.Typ]

	protoZettel      zettel.ProtoZettel
	konfigAkteFormat akten.Format[erworben.Akte, *erworben.Akte]

	Logger
}

type Logger struct {
	New, Updated, Unchanged schnittstellen.FuncIter[*sku.Transacted]
}

func (c *Store) Initialize(
	flags *flag.FlagSet,
	k *konfig.Compiled,
	st standort.Standort,
	pmf objekte_format.Format,
	t thyme.Time,
	virtualStores map[string]*query.VirtualStoreInitable,
	luaVMPoolBuilder *lua.VMPoolBuilder,
) (err error) {
	c.konfig = k
	c.standort = st
	c.akten = akten.Make(st)
	c.persistentMetadateiFormat = pmf
	c.options = objekte_format.Options{Tai: true}
	c.sonnenaufgang = t
	c.fileEncoder = store_fs.MakeFileEncoder(st, k)
	c.virtualStores = virtualStores
	c.luaVMPoolBuilder = luaVMPoolBuilder

	if c.cwdFiles, err = store_fs.MakeCwdFilesAll(
		k,
		sku.StoreFuncs{
			FuncRealize:     c.tryRealize,
			FuncCommit:      c.tryRealizeAndOrStore,
			FuncReadSha:     c.ReadOneEnnui,
			FuncReadOneInto: c.ReadOneInto,
		},
		k.FileExtensions,
		st,
		c.options,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	c.externalStores = map[string]sku.ExternalStore{
		"": c.cwdFiles,
	}

	for k, v := range virtualStores {
		c.externalStores[k] = v
	}

	// c.virtualStores[""] = &query.VirtualStoreInitable{
	// 	VirtualStore: c.cwdFiles,
	// }

	c.metadateiTextParser = metadatei.MakeTextParser(
		c.standort,
		nil, // TODO-P1 make akteFormatter
	)

	c.typenIndex = kennung_index.MakeIndex2[kennung.Typ](
		c.standort,
		st.DirVerzeichnisse("TypenIndexV0"),
	)

	c.bestandsaufnahmeAkte = bestandsaufnahme.Akte{
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

	c.konfigAkteFormat = akten.MakeAkteFormat[erworben.Akte, *erworben.Akte](
		akten.MakeTextParserIgnoreTomlErrors[erworben.Akte](
			c.GetStandort(),
		),
		akten.ParsedAkteTomlFormatter[erworben.Akte, *erworben.Akte]{},
		c.GetStandort(),
	)

	return
}

func (s *Store) SetCheckedOutLogWriter(
	zelw schnittstellen.FuncIter[sku.CheckedOutLike],
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
