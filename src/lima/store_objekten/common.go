package store_objekten

import (
	"bytes"
	"io/ioutil"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/age"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/id"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/echo/age_io"
	"github.com/friedenberg/zit/src/echo/konfig"
	"github.com/friedenberg/zit/src/golf/transaktion"
)

type common struct {
	LockSmith   LockSmith
	Age         age.Age
	Konfig      konfig.Konfig
	Standort    standort.Standort
	Transaktion transaktion.Transaktion
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