package store_objekten

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/age"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/echo/standort"
	"github.com/friedenberg/zit/src/foxtrot/id"
	"github.com/friedenberg/zit/src/golf/age_io"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/juliett/konfig_compiled"
)

type common struct {
	LockSmith   LockSmith
	Age         age.Age
	konfig      *konfig_compiled.Compiled
	Standort    standort.Standort
	Transaktion transaktion.Transaktion
	Abbr        *indexAbbr
}

func (s common) Konfig() konfig_compiled.Compiled {
	return *s.konfig
}

func (s common) KonfigPtr() *konfig_compiled.Compiled {
	return s.konfig
}

func (s common) SizeForObjektenSku(
	sk sku.SkuLike,
) (n int64, err error) {
	var p string

	if p, err = s.Standort.DirObjektenGattung(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	var fi os.FileInfo

	if fi, err = os.Stat(id.Path(sk.GetObjekteSha(), p)); err != nil {
		err = errors.Wrap(err)
		return
	}

	n = fi.Size()

	return
}

func (s common) ReadCloserObjektenSku(
	sk sku.SkuLike,
) (rc sha.ReadCloser, err error) {
	var p string

	if p, err = s.Standort.DirObjektenGattung(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	o := age_io.FileReadOptions{
		Age:  s.Age,
		Path: id.Path(sk.GetObjekteSha(), p),
	}

	if rc, err = age_io.NewFileReader(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s common) WriteCloserObjektenGattung(
	g gattung.GattungLike,
) (wc sha.WriteCloser, err error) {
	var p string

	if p, err = s.Standort.DirObjektenGattung(g.GetGattung()); err != nil {
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

func (s common) AkteReader(sh sha.Sha) (r sha.ReadCloser, err error) {
	if sh.IsNull() {
		r = sha.MakeNopReadCloser(ioutil.NopCloser(bytes.NewReader(nil)))
		return
	}

	p := id.Path(sh, s.Standort.DirObjektenAkten())

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
