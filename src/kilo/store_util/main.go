package store_util

import (
	"bytes"
	"io/ioutil"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/id"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/age"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/delta/age_io"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/kennung_index"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/india/bestandsaufnahme"
	"github.com/friedenberg/zit/src/india/konfig"
)

type StoreUtilVerzeichnisse interface {
	standort.Getter
	konfig.Getter
	schnittstellen.VerzeichnisseFactory
}

type StoreUtil interface {
	errors.Flusher
	StoreUtilVerzeichnisse
	schnittstellen.LockSmithGetter
	konfig.PtrGetter
	schnittstellen.AkteIOFactory
	kennung.Clock

	CommitTransacted(objekte.TransactedLike) error
	CommitUpdatedTransacted(objekte.TransactedLikePtr) error

	GetBestandsaufnahmeStore() bestandsaufnahme.Store
	GetBestandsaufnahmeAkte() bestandsaufnahme.Akte
	GetTransaktionStore() TransaktionStore
	GetAbbrStore() AbbrStore
	GetKennungIndex() kennung_index.Index
	GetEtikettenIndex() (kennung_index.Index2[kennung.Etikett], error)
	GetTypenIndex() (kennung_index.Index2[kennung.Typ], error)

	SetMatchableAdder(kennung.MatchableAdder)
	kennung.MatchableAdder

	objekte_format.Getter

	ObjekteReaderWriterFactory(
		schnittstellen.GattungGetter,
	) schnittstellen.ObjekteIOFactory
}

// TODO-P3 move to own package
type common struct {
	LockSmith                 schnittstellen.LockSmith
	Age                       age.Age
	konfig                    *konfig.Compiled
	standort                  standort.Standort
	transaktion               transaktion.Transaktion
	bestandsaufnahmeAkte      bestandsaufnahme.Akte
	Abbr                      AbbrStore
	persistentMetadateiFormat objekte_format.Format

	bestandsaufnahmeStore bestandsaufnahme.Store
	kennungIndex          kennung_index.Index

	kennung.MatchableAdder
	etikettenIndex verzeichnisseWrapper[kennung_index.Index2[kennung.Etikett]]
	typenIndex     verzeichnisseWrapper[kennung_index.Index2[kennung.Typ]]
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
		etikettenIndex: makeVerzeichnisseWrapper[kennung_index.Index2[kennung.Etikett]](
			kennung_index.MakeIndex2[kennung.Etikett](),
			st.DirVerzeichnisse("EtikettenIndexV0"),
		),
		typenIndex: makeVerzeichnisseWrapper[kennung_index.Index2[kennung.Typ]](
			kennung_index.MakeIndex2[kennung.Typ](),
			st.DirVerzeichnisse("TypenIndexV0"),
		),
	}

	t := kennung.Now()

	for {
		p := c.GetTransaktionStore().TransaktionPath(t)

		if !files.Exists(p) {
			break
		}

		t.MoveForwardIota()
	}

	c.transaktion = transaktion.MakeTransaktion(t)
	c.bestandsaufnahmeAkte = bestandsaufnahme.Akte{
		Skus: sku.MakeSku2Heap(),
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
		c.konfig.GetStoreVersion(),
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

func (s *common) CommitUpdatedTransacted(t objekte.TransactedLikePtr) (err error) {
	ta := kennung.NowTai()
	t.SetTai(ta)

	return s.CommitTransacted(t)
}

func (s *common) CommitTransacted(t objekte.TransactedLike) (err error) {
	s.bestandsaufnahmeAkte.Skus.Add(t.GetSku())

	return
}

func (s *common) GetBestandsaufnahmeAkte() bestandsaufnahme.Akte {
	return s.bestandsaufnahmeAkte
}

func (s *common) GetTransaktionStore() TransaktionStore {
	return s
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

func (s *common) GetEtikettenIndex() (kennung_index.Index2[kennung.Etikett], error) {
	return s.etikettenIndex.Get(s)
}

func (s *common) GetTypenIndex() (kennung_index.Index2[kennung.Typ], error) {
	return s.typenIndex.Get(s)
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

func (s *common) SetMatchableAdder(ma kennung.MatchableAdder) {
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
		Age:  s.Age,
		Path: id.Path(sh.GetSha(), p),
	}

	if rc, err = age_io.NewFileReader(o); err != nil {
		err = errors.Wrapf(err, "Gattung: %s", g.GetGattung())
		err = errors.Wrapf(err, "Sha: %s", sh.GetSha())
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
		LockFile:                 true,
	}

	if wc, err = age_io.NewMover(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s common) ReadCloserObjekten(p string) (sha.ReadCloser, error) {
	o := age_io.FileReadOptions{
		Age:  s.Age,
		Path: p,
	}

	return age_io.NewFileReader(o)
}

func (s common) ReadCloserVerzeichnisse(p string) (sha.ReadCloser, error) {
	o := age_io.FileReadOptions{
		Age:  s.Age,
		Path: p,
	}

	return age_io.NewFileReader(o)
}

func (s common) WriteCloserObjekten(p string) (w sha.WriteCloser, err error) {
	return age_io.NewMover(
		age_io.MoveOptions{
			Age:       s.Age,
			FinalPath: p,
			LockFile:  true,
		},
	)
}

func (s common) WriteCloserVerzeichnisse(p string) (w sha.WriteCloser, err error) {
	return age_io.NewMover(
		age_io.MoveOptions{
			Age:       s.Age,
			FinalPath: p,
			LockFile:  false,
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
		LockFile:                 true,
	}

	if outer, err = age_io.NewMover(mo); err != nil {
		err = errors.Wrap(err)
		return
	}

	w = outer

	return
}

func (s common) AkteReader(sh sha.ShaLike) (r sha.ReadCloser, err error) {
	if sh.GetSha().IsNull() {
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

	p = id.Path(sh.GetSha(), p)

	o := age_io.FileReadOptions{
		Age:  s.Age,
		Path: p,
	}

	if r, err = age_io.NewFileReader(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
