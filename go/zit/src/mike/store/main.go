package store

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/delta/standort"
	"code.linenisgreat.com/zit/src/delta/thyme"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/src/golf/kennung_index"
	"code.linenisgreat.com/zit/src/golf/objekte_format"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/erworben"
	"code.linenisgreat.com/zit/src/india/objekte_collections"
	"code.linenisgreat.com/zit/src/india/query"
	"code.linenisgreat.com/zit/src/juliett/konfig"
	"code.linenisgreat.com/zit/src/juliett/objekte"
	"code.linenisgreat.com/zit/src/kilo/cwd"
	"code.linenisgreat.com/zit/src/kilo/objekte_store"
	"code.linenisgreat.com/zit/src/kilo/store_verzeichnisse"
	"code.linenisgreat.com/zit/src/kilo/zettel"
	"code.linenisgreat.com/zit/src/lima/akten"
	"code.linenisgreat.com/zit/src/lima/bestandsaufnahme"
)

type Store struct {
	konfig                    *konfig.Compiled
	standort                  standort.Standort
	cwdFiles                  *cwd.CwdFiles
	akten                     *akten.Akten
	bestandsaufnahmeAkte      bestandsaufnahme.Akte
	options                   objekte_format.Options
	Abbr                      AbbrStore
	persistentMetadateiFormat objekte_format.Format
	fileEncoder               objekte_collections.FileEncoder

	verzeichnisse *store_verzeichnisse.Store

	sonnenaufgang thyme.Time

	checkedOutLogPrinter schnittstellen.FuncIter[*sku.CheckedOut]

	metadateiTextParser metadatei.TextParser

	bestandsaufnahmeStore bestandsaufnahme.Store
	kennungIndex          kennung_index.Index

	query.MatchableAdder
	typenIndex kennung_index.KennungIndex[kennung.Typ, *kennung.Typ]

	protoZettel      zettel.ProtoZettel
	konfigAkteFormat objekte.AkteFormat[erworben.Akte, *erworben.Akte]

	objekte_store.LogWriter
}

func (c *Store) Initialize(
	k *konfig.Compiled,
	st standort.Standort,
	pmf objekte_format.Format,
	t thyme.Time,
) (err error) {
	c.konfig = k
	c.standort = st
	c.akten = akten.Make(st)
	c.persistentMetadateiFormat = pmf
	c.options = objekte_format.Options{Tai: true}
	c.sonnenaufgang = t
	c.fileEncoder = objekte_collections.MakeFileEncoder(st, k)

	if c.cwdFiles, err = cwd.MakeCwdFilesAll(
		k,
		st,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

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
		c.GetKennungIndex(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	c.protoZettel = zettel.MakeProtoZettel(c.GetKonfig())

	c.konfigAkteFormat = objekte_store.MakeAkteFormat[erworben.Akte, *erworben.Akte](
		objekte.MakeTextParserIgnoreTomlErrors[erworben.Akte](
			c.GetStandort(),
		),
		objekte.ParsedAkteTomlFormatter[erworben.Akte, *erworben.Akte]{},
		c.GetStandort(),
	)

	return
}

func (s *Store) SetCheckedOutLogWriter(
	zelw schnittstellen.FuncIter[*sku.CheckedOut],
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

func (s *Store) SetLogWriter(lw objekte_store.LogWriter) {
	s.LogWriter = lw
}
