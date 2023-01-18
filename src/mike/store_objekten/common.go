package store_objekten

import (
	"bytes"
	"io/ioutil"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/age"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/delta/id"
	"github.com/friedenberg/zit/src/echo/age_io"
	"github.com/friedenberg/zit/src/foxtrot/standort"
	"github.com/friedenberg/zit/src/foxtrot/ts"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/india/bestandsaufnahme"
	"github.com/friedenberg/zit/src/juliett/konfig"
)

// TODO-P3 move to own package
type common struct {
	LockSmith        LockSmith
	Age              age.Age
	konfig           *konfig.Compiled
	Standort         standort.Standort
	Transaktion      transaktion.Transaktion
	Bestandsaufnahme *bestandsaufnahme.Objekte
	Abbr             *indexAbbr

	bestandsaufnahmeStore bestandsaufnahme.Store
}

func (s common) AddSkuToBestandsaufnahme(sk sku.SkuLike, as sha.Sha) {
	s.Bestandsaufnahme.Akte.Skus.Push(
		sku.Sku2{
			Gattung:    gattung.Make(sk.GetGattung()),
			Tai:        ts.TaiFromTimeWithIndex(sk.GetTime(), sk.GetTransactionIndex().Int()),
			Kennung:    collections.MakeStringValue(sk.GetId().String()),
			ObjekteSha: sha.Make(sk.GetObjekteSha()),
			AkteSha:    as,
		},
	)
}

func (s common) GetKonfig() konfig.Compiled {
	return *s.konfig
}

func (s common) Konfig() konfig.Compiled {
	return *s.konfig
}

func (s common) KonfigPtr() *konfig.Compiled {
	return s.konfig
}

func (s common) GetKonfigPtr() *konfig.Compiled {
	return s.konfig
}

func (s common) ObjekteReader(
	g gattung.GattungLike,
	sh sha.ShaLike,
) (rc sha.ReadCloser, err error) {
	var p string

	if p, err = s.Standort.DirObjektenGattung(g); err != nil {
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
	g gattung.GattungLike,
) (wc sha.WriteCloser, err error) {
	var p string

	if p, err = s.Standort.DirObjektenGattung(g); err != nil {
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
		FinalPath:                s.Standort.DirObjektenAkten(),
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

	p := id.Path(sh.GetSha(), s.Standort.DirObjektenAkten())

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
