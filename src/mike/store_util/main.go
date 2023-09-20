package store_util

import (
	"bytes"
	"io/ioutil"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/id"
	"github.com/friedenberg/zit/src/charlie/age"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/pool"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/echo/age_io"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/golf/kennung_index"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/kilo/bestandsaufnahme"
	"github.com/friedenberg/zit/src/kilo/konfig"
)

type StoreUtilVerzeichnisse interface {
	GetStoreVersion() schnittstellen.StoreVersion
	standort.Getter
	konfig.Getter
	schnittstellen.VerzeichnisseFactory
}

type StoreUtil interface {
	GetSkuPool() schnittstellen.Pool[sku.Transacted2, *sku.Transacted2]
	FlushBestandsaufnahme() error
	errors.Flusher
	StoreUtilVerzeichnisse
	schnittstellen.LockSmithGetter
	konfig.PtrGetter
	schnittstellen.AkteIOFactory
	kennung.Clock

	ExternalReader
	CommitTransacted(sku.SkuLike) error
	CommitUpdatedTransacted(sku.SkuLikePtr) error

	GetBestandsaufnahmeStore() bestandsaufnahme.Store
	GetBestandsaufnahmeAkte() bestandsaufnahme.Akte
	GetAbbrStore() AbbrStore
	GetKennungIndex() kennung_index.Index
	GetTypenIndex() (kennung_index.KennungIndex[kennung.Typ, *kennung.Typ], error)

	SetMatchableAdder(matcher.MatchableAdder)
	matcher.MatchableAdder

	objekte_format.Getter

	ObjekteReaderWriterFactory(
		schnittstellen.GattungGetter,
	) schnittstellen.ObjekteIOFactory

	GetExternalReader2() *ExternalReader2
}

// TODO-P3 move to own package
type common struct {
	LockSmith                 schnittstellen.LockSmith
	Age                       age.Age
	konfig                    *konfig.Compiled
	standort                  standort.Standort
	bestandsaufnahmeAkte      bestandsaufnahme.Akte
	Abbr                      AbbrStore
	persistentMetadateiFormat objekte_format.Format
	pool                      schnittstellen.Pool[sku.Transacted2, *sku.Transacted2]

	metadateiTextParser metadatei.TextParser

	bestandsaufnahmeStore bestandsaufnahme.Store
	kennungIndex          kennung_index.Index

	matcher.MatchableAdder
	typenIndex kennung_index.KennungIndex[kennung.Typ, *kennung.Typ]

	er2 ExternalReader2
}

func MakeStoreUtil(
	lockSmith schnittstellen.LockSmith,
	a age.Age,
	k *konfig.Compiled,
	st standort.Standort,
	pmf objekte_format.Format,
) (c *common, err error) {
	c = &common{
		LockSmith:                 lockSmith,
		Age:                       a,
		konfig:                    k,
		standort:                  st,
		persistentMetadateiFormat: pmf,
		pool: pool.MakePool[sku.Transacted2, *sku.Transacted2](
			nil,
			nil,
		),
	}

	c.metadateiTextParser = metadatei.MakeTextParser(
		c,
		nil, // TODO-P1 make akteFormatter
	)

	c.er2 = ExternalReader2{
		metadateiTextParser: c.metadateiTextParser,
		AkteIOFactory:       c,
	}

	c.typenIndex = kennung_index.MakeIndex2[kennung.Typ](
		c,
		st.DirVerzeichnisse("TypenIndexV0"),
	)

	c.bestandsaufnahmeAkte = bestandsaufnahme.Akte{
		Skus: sku.MakeSkuLikeHeap(),
	}

	if c.Abbr, err = newIndexAbbr(
		c,
		st.DirVerzeichnisse("Abbr"),
	); err != nil {
		err = errors.Wrapf(err, "failed to init abbr index")
		return
	}

	if c.bestandsaufnahmeStore, err = bestandsaufnahme.MakeStore(
		c.GetStandort(),
		c.GetLockSmith(),
		c.konfig.GetStoreVersion(),
		c,
		c.ObjekteReaderWriterFactory(gattung.Bestandsaufnahme),
		c,
		pmf,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if c.kennungIndex, err = kennung_index.MakeIndex(
		c.GetKonfig(),
		c.GetStandort(),
		c,
	); err != nil {
		err = errors.Wrapf(err, "failed to init zettel index")
		return
	}

	return
}

func (s *common) GetSkuPool() schnittstellen.Pool[sku.Transacted2, *sku.Transacted2] {
	return s.pool
}

func (s *common) GetExternalReader2() *ExternalReader2 {
	return &s.er2
}

func (s common) GetStoreVersion() schnittstellen.StoreVersion {
	return s.konfig.StoreVersion
}

func (s common) GetLockSmith() schnittstellen.LockSmith {
	return s.LockSmith
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
	sk := t.GetSkuLike()
	sku.AddSkuToHeap(&s.bestandsaufnahmeAkte.Skus, sk)

	return
}

func (s *common) GetBestandsaufnahmeAkte() bestandsaufnahme.Akte {
	return s.bestandsaufnahmeAkte
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

func (s common) objekteReader(
	g schnittstellen.GattungGetter,
	sh sha.ShaLike,
) (rc sha.ReadCloser, err error) {
	var p string

	if p, err = s.GetStandort().DirObjektenGattung(
		s.konfig.GetStoreVersion(),
		g,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	o := age_io.FileReadOptions{
		Age:             s.Age,
		Path:            id.Path(sh.GetShaLike(), p),
		CompressionType: s.GetKonfig().CompressionType,
	}

	if rc, err = age_io.NewFileReader(o); err != nil {
		err = errors.Wrapf(err, "Gattung: %s", g.GetGattung())
		err = errors.Wrapf(err, "Sha: %s", sh.GetShaLike())
		return
	}

	return
}

func (s common) objekteWriter(
	g schnittstellen.GattungGetter,
) (wc sha.WriteCloser, err error) {
	var p string

	if p, err = s.GetStandort().DirObjektenGattung(
		s.konfig.GetStoreVersion(),
		g,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	o := age_io.MoveOptions{
		Age:                      s.Age,
		FinalPath:                p,
		GenerateFinalPathFromSha: true,
		LockFile:                 s.GetKonfig().LockInternalFiles,
		CompressionType:          s.GetKonfig().CompressionType,
	}

	if wc, err = age_io.NewMover(s.GetStandort(), o); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s common) ReadCloserObjekten(p string) (sha.ReadCloser, error) {
	o := age_io.FileReadOptions{
		Age:             s.Age,
		Path:            p,
		CompressionType: s.GetKonfig().CompressionType,
	}

	return age_io.NewFileReader(o)
}

func (s common) ReadCloserVerzeichnisse(p string) (sha.ReadCloser, error) {
	o := age_io.FileReadOptions{
		Age:             s.Age,
		Path:            p,
		CompressionType: s.GetKonfig().CompressionType,
	}

	return age_io.NewFileReader(o)
}

func (s common) WriteCloserObjekten(p string) (w sha.WriteCloser, err error) {
	return age_io.NewMover(
		s.GetStandort(),
		age_io.MoveOptions{
			Age:             s.Age,
			FinalPath:       p,
			LockFile:        s.GetKonfig().LockInternalFiles,
			CompressionType: s.GetKonfig().CompressionType,
		},
	)
}

func (s common) WriteCloserVerzeichnisse(
	p string,
) (w sha.WriteCloser, err error) {
	return age_io.NewMover(
		s.GetStandort(),
		age_io.MoveOptions{
			Age:             s.Age,
			FinalPath:       p,
			LockFile:        false,
			CompressionType: s.GetKonfig().CompressionType,
		},
	)
}

func (s common) AkteWriter() (w sha.WriteCloser, err error) {
	var outer age_io.Writer

	var p string

	if p, err = s.standort.DirObjektenGattung(
		s.konfig.GetStoreVersion(),
		gattung.Akte,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	mo := age_io.MoveOptions{
		Age:                      s.Age,
		FinalPath:                p,
		GenerateFinalPathFromSha: true,
		LockFile:                 s.GetKonfig().LockInternalFiles,
		CompressionType:          s.GetKonfig().CompressionType,
	}

	if outer, err = age_io.NewMover(s.GetStandort(), mo); err != nil {
		err = errors.Wrap(err)
		return
	}

	w = outer

	return
}

func (s common) AkteReader(sh sha.ShaLike) (r sha.ReadCloser, err error) {
	if sh.GetShaLike().IsNull() {
		r = sha.MakeNopReadCloser(ioutil.NopCloser(bytes.NewReader(nil)))
		return
	}

	var p string

	if p, err = s.standort.DirObjektenGattung(
		s.konfig.GetStoreVersion(),
		gattung.Akte,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	p = id.Path(sh.GetShaLike(), p)

	o := age_io.FileReadOptions{
		Age:             s.Age,
		Path:            p,
		CompressionType: s.GetKonfig().CompressionType,
	}

	if r, err = age_io.NewFileReader(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
