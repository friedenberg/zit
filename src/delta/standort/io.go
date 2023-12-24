package standort

import (
	"bytes"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/id"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
)

func (s Standort) objekteReader(
	g schnittstellen.GattungGetter,
	sh sha.ShaLike,
) (rc sha.ReadCloser, err error) {
	var p string

	if p, err = s.DirObjektenGattung(
		s.angeboren.GetStoreVersion(),
		g,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	o := FileReadOptions{
		Age:             *s.age,
		Path:            id.Path(sh.GetShaLike(), p),
		CompressionType: s.angeboren.CompressionType,
	}

	if rc, err = NewFileReader(o); err != nil {
		err = errors.Wrapf(err, "Gattung: %s", g.GetGattung())
		err = errors.Wrapf(err, "Sha: %s", sh.GetShaLike())
		return
	}

	return
}

func (s Standort) objekteWriter(
	g schnittstellen.GattungGetter,
) (wc sha.WriteCloser, err error) {
	var p string

	if p, err = s.DirObjektenGattung(
		s.angeboren.GetStoreVersion(),
		g,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	o := MoveOptions{
		Age:                      *s.age,
		FinalPath:                p,
		GenerateFinalPathFromSha: true,
		LockFile:                 s.angeboren.LockInternalFiles,
		CompressionType:          s.angeboren.CompressionType,
	}

	if wc, err = s.NewMover(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Standort) ReadCloserObjekten(p string) (sha.ReadCloser, error) {
	o := FileReadOptions{
		Age:             *s.age,
		Path:            p,
		CompressionType: s.angeboren.CompressionType,
	}

	return NewFileReader(o)
}

func (s Standort) ReadCloserVerzeichnisse(p string) (sha.ReadCloser, error) {
	o := FileReadOptions{
		Age:             *s.age,
		Path:            p,
		CompressionType: s.angeboren.CompressionType,
	}

	return NewFileReader(o)
}

func (s Standort) WriteCloserObjekten(p string) (w sha.WriteCloser, err error) {
	return s.NewMover(
		MoveOptions{
			Age:             *s.age,
			FinalPath:       p,
			LockFile:        s.angeboren.LockInternalFiles,
			CompressionType: s.angeboren.CompressionType,
		},
	)
}

func (s Standort) WriteCloserVerzeichnisse(
	p string,
) (w sha.WriteCloser, err error) {
	return s.NewMover(
		MoveOptions{
			Age:             *s.age,
			FinalPath:       p,
			LockFile:        false,
			CompressionType: s.angeboren.CompressionType,
		},
	)
}

func (s Standort) AkteWriter() (w sha.WriteCloser, err error) {
	var outer Writer

	var p string

	if p, err = s.DirObjektenGattung(
		s.angeboren.GetStoreVersion(),
		gattung.Akte,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	mo := MoveOptions{
		Age:                      *s.age,
		FinalPath:                p,
		GenerateFinalPathFromSha: true,
		LockFile:                 s.angeboren.LockInternalFiles,
		CompressionType:          s.angeboren.CompressionType,
	}

	if outer, err = s.NewMover(mo); err != nil {
		err = errors.Wrap(err)
		return
	}

	w = outer

	return
}

func (s Standort) AkteReader(sh sha.ShaLike) (r sha.ReadCloser, err error) {
	if sh.GetShaLike().IsNull() {
		r = sha.MakeNopReadCloser(io.NopCloser(bytes.NewReader(nil)))
		return
	}

	var p string

	if p, err = s.DirObjektenGattung(
		s.angeboren.GetStoreVersion(),
		gattung.Akte,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	p = id.Path(sh.GetShaLike(), p)

	o := FileReadOptions{
		Age:             *s.age,
		Path:            p,
		CompressionType: s.angeboren.CompressionType,
	}

	if r, err = NewFileReader(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
