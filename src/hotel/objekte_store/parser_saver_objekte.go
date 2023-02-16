package objekte_store

import (
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/golf/objekte"
)

type ParseSaver[
	T schnittstellen.Objekte[T],
	T1 schnittstellen.ObjektePtr[T],
] interface {
	ParseAndSaveAkteAndObjekte(string) (T, schnittstellen.Sha, error)
}

type objekteParseSaver[
	T schnittstellen.Objekte[T],
	T1 schnittstellen.ObjektePtr[T],
] struct {
	awf          schnittstellen.AkteWriterFactory
	akteParser   schnittstellen.Parser[T, T1]
	objekteSaver ObjekteSaver[T, T1]
}

func MakeParseSaver[
	T schnittstellen.Objekte[T],
	T1 schnittstellen.ObjektePtr[T],
](
	owf schnittstellen.ObjekteWriterFactory,
	awf schnittstellen.AkteWriterFactory,
	akteParser schnittstellen.Parser[T, T1],
) *objekteParseSaver[T, T1] {
	return &objekteParseSaver[T, T1]{
		awf:        awf,
		akteParser: akteParser,
		objekteSaver: MakeObjekteSaver[T, T1](
			owf,
			objekte.Format[T, T1]{},
		),
	}
}

func (h *objekteParseSaver[T, T1]) ParseAndSaveAkteAndObjekte(
	p string,
) (o T, sh schnittstellen.Sha, err error) {
	var f *os.File

	if f, err = files.OpenExclusiveReadOnly(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	r := sha.MakeReadCloser(f)

	defer errors.DeferredCloser(&err, r)

	if err = h.readAkte(r, T1(&o)); err != nil {
		err = errors.Wrap(err)
		return
	}

	T1(&o).SetAkteSha(r.Sha())

	if sh, err = h.objekteSaver.SaveObjekte(&o); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (h *objekteParseSaver[T, T1]) readAkte(
	r sha.ReadCloser,
	o T1,
) (err error) {
	var n int64

	if n, err = h.akteParser.Parse(r, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Log().Printf("parsed %d akte bytes", n)

	return
}
