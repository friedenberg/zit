package store_util

import (
	"bytes"
	"io/ioutil"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/id"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/age"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/delta/age_io"
	"github.com/friedenberg/zit/src/echo/ts"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/india/bestandsaufnahme"
	"github.com/friedenberg/zit/src/india/konfig"
	"github.com/friedenberg/zit/src/kennung_index"
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
	schnittstellen.ObjekteAkteFactory
	ts.Clock

	CommitTransacted(objekte.TransactedLike) error

	GetBestandsaufnahmeStore() bestandsaufnahme.Store
	GetBestandsaufnahme() *bestandsaufnahme.Objekte
	GetTransaktionStore() TransaktionStore
	GetAbbrStore() AbbrStore
	GetKennungIndex() kennung_index.Index
}

// TODO-P3 move to own package
type common struct {
	LockSmith        schnittstellen.LockSmith
	Age              age.Age
	konfig           *konfig.Compiled
	standort         standort.Standort
	transaktion      transaktion.Transaktion
	bestandsaufnahme *bestandsaufnahme.Objekte
	Abbr             *indexAbbr

	bestandsaufnahmeStore bestandsaufnahme.Store
	kennungIndex          kennung_index.Index
}

func MakeStoreUtil(
	lockSmith schnittstellen.LockSmith,
	a age.Age,
	k *konfig.Compiled,
	st standort.Standort,
) (c *common, err error) {
	c = &common{
		LockSmith: lockSmith,
		Age:       a,
		konfig:    k,
		standort:  st,
	}

	t := ts.Now()
	ta := ts.NowTai()

	for {
		p := c.GetTransaktionStore().TransaktionPath(t)

		if !files.Exists(p) {
			break
		}

		t.MoveForwardIota()
	}

	c.transaktion = transaktion.MakeTransaktion(t)
	c.bestandsaufnahme = &bestandsaufnahme.Objekte{
		Tai: ta,
		Akte: bestandsaufnahme.Akte{
			Skus: sku.MakeSku2Heap(),
		},
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
		c,
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

func (s common) GetTime() ts.Time {
	return s.transaktion.Time
}

func (s common) GetTai() ts.Tai {
	return s.bestandsaufnahme.Tai
}

func (s common) CommitTransacted(t objekte.TransactedLike) (err error) {
	s.GetBestandsaufnahme().Akte.Skus.Add(t.GetSku2())
	s.GetTransaktionStore().GetTransaktion().Skus.Add(t.GetSkuLike())

	return
}

func (s *common) GetBestandsaufnahme() *bestandsaufnahme.Objekte {
	return s.bestandsaufnahme
}

func (s *common) GetTransaktion() *transaktion.Transaktion {
	return &s.transaktion
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

func (s common) GetStandort() standort.Standort {
	return s.standort
}

func (s common) GetKonfig() konfig.Compiled {
	return *s.konfig
}

func (s common) GetKonfigPtr() *konfig.Compiled {
	return s.konfig
}

func (s common) ObjekteReader(
	g schnittstellen.GattungGetter,
	sh sha.ShaLike,
) (rc sha.ReadCloser, err error) {
	var p string

	if p, err = s.GetStandort().DirObjektenGattung(g); err != nil {
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

func (s common) ObjekteWriter(
	g schnittstellen.GattungGetter,
) (wc sha.WriteCloser, err error) {
	var p string

	if p, err = s.GetStandort().DirObjektenGattung(g); err != nil {
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

	mo := age_io.MoveOptions{
		Age:                      s.Age,
		FinalPath:                s.GetStandort().DirObjektenAkten(),
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

	p := id.Path(sh.GetSha(), s.GetStandort().DirObjektenAkten())

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
