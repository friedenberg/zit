package store_util

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/delta/standort"
	"code.linenisgreat.com/zit/src/delta/thyme"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/src/golf/kennung_index"
	"code.linenisgreat.com/zit/src/golf/objekte_format"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/matcher"
	"code.linenisgreat.com/zit/src/india/objekte_collections"
	"code.linenisgreat.com/zit/src/juliett/konfig"
	"code.linenisgreat.com/zit/src/kilo/cwd"
	"code.linenisgreat.com/zit/src/kilo/store_verzeichnisse"
	"code.linenisgreat.com/zit/src/lima/akten"
	"code.linenisgreat.com/zit/src/lima/bestandsaufnahme"
)

type StoreUtil interface {
	FlushBestandsaufnahme() error
	errors.FlusherWithLogger
	kennung.Clock

	ExternalReader

	mutators
	accessors
	reader

	ResetIndexes() (err error)
	AddTypToIndex(t *kennung.Typ) (err error)

	SetMatchableAdder(matcher.MatchableAdder)
	matcher.MatchableAdder

	SetCheckedOutLogWriter(zelw schnittstellen.FuncIter[*sku.CheckedOut])

	ReadOneExternalFS(*sku.Transacted) (*sku.CheckedOut, error)

	CheckoutQuery(
		options checkout_options.Options,
		fq matcher.FuncReaderTransactedLikePtr,
		f schnittstellen.FuncIter[*sku.CheckedOut],
	) (err error)

	Checkout(
		options checkout_options.Options,
		fq matcher.FuncReaderTransactedLikePtr,
		ztw schnittstellen.FuncIter[*sku.Transacted],
	) (zcs sku.CheckedOutMutableSet, err error)

	ReadFiles(
		fq matcher.FuncReaderTransactedLikePtr,
		f schnittstellen.FuncIter[*sku.CheckedOut],
	) (err error)

	CheckoutOne(
		options checkout_options.Options,
		sz *sku.Transacted,
	) (cz *sku.CheckedOut, err error)
}

// TODO-P3 move to own package
type common struct {
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

	matcher.MatchableAdder
	typenIndex kennung_index.KennungIndex[kennung.Typ, *kennung.Typ]
}

func MakeStoreUtil(
	k *konfig.Compiled,
	st standort.Standort,
	pmf objekte_format.Format,
	t thyme.Time,
) (c *common, err error) {
	c = &common{
		konfig:                    k,
		standort:                  st,
		akten:                     akten.Make(st),
		persistentMetadateiFormat: pmf,
		options:                   objekte_format.Options{Tai: true},
		sonnenaufgang:             t,
		fileEncoder:               objekte_collections.MakeFileEncoder(st, k),
	}

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

	return
}

func (s *common) SetCheckedOutLogWriter(
	zelw schnittstellen.FuncIter[*sku.CheckedOut],
) {
	s.checkedOutLogPrinter = zelw
}

func (s *common) ResetIndexes() (err error) {
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

func (s *common) SetMatchableAdder(ma matcher.MatchableAdder) {
	s.MatchableAdder = ma
}
