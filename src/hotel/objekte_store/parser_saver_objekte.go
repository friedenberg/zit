package objekte_store

import (
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
)

type ParseSaver[
	T objekte.Objekte[T],
	T1 objekte.ObjektePtr[T],
	T2 schnittstellen.Id[T2],
	T3 schnittstellen.IdPtr[T2],
] interface {
	ParseAndSaveAkteAndObjekte(sku.ExternalMaybe[T2, T3]) (T, sku.External[T2, T3], error)
}

type objekteParseSaver[
	T objekte.Objekte[T],
	T1 objekte.ObjektePtr[T],
	T2 schnittstellen.Id[T2],
	T3 schnittstellen.IdPtr[T2],
] struct {
	awf          schnittstellen.AkteWriterFactory
	akteParser   schnittstellen.Parser[T, T1]
	objekteSaver ObjekteSaver[T, T1]
}

func MakeParseSaver[
	T objekte.Objekte[T],
	T1 objekte.ObjektePtr[T],
	T2 schnittstellen.Id[T2],
	T3 schnittstellen.IdPtr[T2],
](
	owf schnittstellen.ObjekteIOFactory,
	awf schnittstellen.AkteIOFactory,
	akteParser schnittstellen.Parser[T, T1],
) *objekteParseSaver[T, T1, T2, T3] {
	if akteParser == nil {
		akteParser = MakeNopAkteFormat[T, T1](awf)
	}

	return &objekteParseSaver[T, T1, T2, T3]{
		awf:        awf,
		akteParser: akteParser,
		objekteSaver: MakeObjekteSaver[T, T1](
			owf,
			objekte.Format[T, T1]{},
		),
	}
}

func (h *objekteParseSaver[T, T1, T2, T3]) ParseAndSaveAkteAndObjekte(
	sem sku.ExternalMaybe[T2, T3],
) (o T, sk sku.External[T2, T3], err error) {
	var f *os.File

	errors.TodoP2("support akte")
	if f, err = files.OpenExclusiveReadOnly(sem.FDs.Objekte.Path); err != nil {
		err = errors.Wrap(err)
		return
	}

	sk.FDs = sem.FDs
	sk.Kennung = sem.Kennung

	errors.TodoP0("populate sku.AkteSha and ObjekteSha")

	r := sha.MakeReadCloser(f)

	defer errors.DeferredCloser(&err, r)

	if err = h.readAkte(r, T1(&o)); err != nil {
		err = errors.Wrap(err)
		return
	}

	T1(&o).SetAkteSha(r.Sha())

	var sh schnittstellen.Sha

	if sh, err = h.objekteSaver.SaveObjekte(&o); err != nil {
		err = errors.Wrap(err)
		return
	}

	sk.ObjekteSha = sha.Make(sh)

	return
}

func (h *objekteParseSaver[T, T1, T2, T3]) readAkte(
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
