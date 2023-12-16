package store_util

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/checkout_options"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/delta/thyme"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/golf/kennung_index"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/india/objekte_collections"
	"github.com/friedenberg/zit/src/juliett/konfig"
	"github.com/friedenberg/zit/src/kilo/cwd"
	"github.com/friedenberg/zit/src/kilo/store_verzeichnisse"
	"github.com/friedenberg/zit/src/lima/akten"
	"github.com/friedenberg/zit/src/lima/bestandsaufnahme"
)

type StoreUtil interface {
	FlushBestandsaufnahme() error
	errors.Flusher
	kennung.Clock

	ExternalReader
	CommitTransacted(*sku.Transacted) error
	CommitUpdatedTransacted(*sku.Transacted) error
	CalculateAndSetShaTransacted(sk *sku.Transacted) (err error)
	CalculateAndSetShaSkuLike(sk sku.SkuLike) (err error)

	accessors

	ResetIndexes() (err error)
	AddTypToIndex(t *kennung.Typ) (err error)

	ReadAllGattung(
		g gattung.Gattung,
		f schnittstellen.FuncIter[*sku.Transacted],
	) (err error)

	ReadAllGattungen(
		g gattungen.Set,
		f schnittstellen.FuncIter[*sku.Transacted],
	) (err error)

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

	verzeichnisseSchwanzen *VerzeichnisseSchwanzen
	verzeichnisseAll       *store_verzeichnisse.Store

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
		options:                   objekte_format.Options{IncludeTai: true},
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
		c.standort,
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

	return
}

func (s *common) SetCheckedOutLogWriter(
	zelw schnittstellen.FuncIter[*sku.CheckedOut],
) {
	s.checkedOutLogPrinter = zelw
}

func (s *common) CommitUpdatedTransacted(
	t *sku.Transacted,
) (err error) {
	ta := kennung.NowTai()
	t.SetTai(ta)

	return s.CommitTransacted(t)
}

func (s *common) CommitTransacted(t *sku.Transacted) (err error) {
	sk := sku.GetTransactedPool().Get()

	if err = sk.SetFromSkuLike(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.bestandsaufnahmeAkte.Skus.Add(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
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

func (s *common) AddTypToIndex(t *kennung.Typ) (err error) {
	if t == nil {
		return
	}

	if t.IsEmpty() {
		return
	}

	if err = s.typenIndex.StoreOne(*t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *common) SetMatchableAdder(ma matcher.MatchableAdder) {
	s.MatchableAdder = ma
}
