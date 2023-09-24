package store_util

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/delta/typ_akte"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/golf/kennung_index"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/kilo/konfig"
	"github.com/friedenberg/zit/src/lima/cwd"
	"github.com/friedenberg/zit/src/mike/bestandsaufnahme"
)

type StoreUtil interface {
	FlushBestandsaufnahme() error
	errors.Flusher
	standort.Getter
	konfig.Getter
	konfig.PtrGetter
	kennung.Clock

	ExternalReader
	CommitTransacted(sku.SkuLike) error
	CommitUpdatedTransacted(sku.SkuLikePtr) error

	GetBestandsaufnahmeStore() bestandsaufnahme.Store
	GetAbbrStore() AbbrStore
	GetKennungIndex() kennung_index.Index
	GetTypenIndex() (kennung_index.KennungIndex[kennung.Typ, *kennung.Typ], error)

	SetMatchableAdder(matcher.MatchableAdder)
	matcher.MatchableAdder

	objekte_format.Getter

	SetCheckedOutLogWriter(zelw schnittstellen.FuncIter[*sku.CheckedOut])

	CheckoutQuery(
		options CheckoutOptions,
		fq matcher.FuncReaderTransactedLikePtr,
		f schnittstellen.FuncIter[*sku.CheckedOut],
	) (err error)

	Checkout(
		options CheckoutOptions,
		tagp schnittstellen.AkteGetterPutter[*typ_akte.V0],
		fq matcher.FuncReaderTransactedLikePtr,
		ztw schnittstellen.FuncIter[*sku.Transacted],
	) (zcs schnittstellen.MutableSetLike[*sku.CheckedOut], err error)

	ReadFiles(
		fs *cwd.CwdFiles,
		fq matcher.FuncReaderTransactedLikePtr,
		f schnittstellen.FuncIter[*sku.CheckedOut],
	) (err error)

	CheckoutOne(
		options CheckoutOptions,
		sz *sku.Transacted,
	) (cz *sku.CheckedOut, err error)
}

// TODO-P3 move to own package
type common struct {
	konfig                    *konfig.Compiled
	standort                  standort.Standort
	bestandsaufnahmeAkte      bestandsaufnahme.Akte
	Abbr                      AbbrStore
	persistentMetadateiFormat objekte_format.Format

	sonnenaufgang kennung.Time

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
	t kennung.Time,
) (c *common, err error) {
	c = &common{
		konfig:                    k,
		standort:                  st,
		persistentMetadateiFormat: pmf,
		sonnenaufgang:             t,
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

func (s common) GetPersistentMetadateiFormat() objekte_format.Format {
	return s.persistentMetadateiFormat
}

func (s common) GetTime() kennung.Time {
	return kennung.Now()
}

func (s common) GetTai() kennung.Tai {
	return kennung.NowTai()
}

func (s *common) CommitUpdatedTransacted(
	t sku.SkuLikePtr,
) (err error) {
	ta := kennung.NowTai()
	t.SetTai(ta)

	return s.CommitTransacted(t)
}

func (s *common) CommitTransacted(t sku.SkuLike) (err error) {
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

func (s *common) GetBestandsaufnahmeStore() bestandsaufnahme.Store {
	return s.bestandsaufnahmeStore
}

func (s *common) GetAbbrStore() AbbrStore {
	return s.Abbr
}

func (s *common) GetKennungIndex() kennung_index.Index {
	return s.kennungIndex
}

func (s *common) GetTypenIndex() (kennung_index.KennungIndex[kennung.Typ, *kennung.Typ], error) {
	return s.typenIndex, nil
}

func (s common) GetStandort() standort.Standort {
	return s.standort
}

func (s common) GetKonfig() konfig.Compiled {
	return *s.konfig
}

func (s common) GetKonfigPtr() *konfig.Compiled {
	return s.konfig
}

func (s *common) SetMatchableAdder(ma matcher.MatchableAdder) {
	s.MatchableAdder = ma
}
