package bestandsaufnahme

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/echo/standort"
	"github.com/friedenberg/zit/src/hotel/objekte"
)

type Store interface {
	ObjekteInflator
	ObjekteSaver
	AkteTextSaver
	Create(*Objekte) (schnittstellen.Sha, error)
	ReadOne(sha.Sha) (*Objekte, error)
	ReadAll(collections.WriterFunc[*Objekte]) error
	// errors.Flusher
}

type ObjekteSaver = objekte.ObjekteSaver[
	Objekte,
	*Objekte,
]

type AkteTextSaver = objekte.AkteTextSaver[
	Objekte,
	*Objekte,
]

type ObjekteInflator = objekte.ObjekteInflator[
	Objekte,
	*Objekte,
	objekte.NilVerzeichnisse[Objekte],
	*objekte.NilVerzeichnisse[Objekte],
]

type ObjekteFormat = schnittstellen.Format[
	Objekte,
	*Objekte,
]

type AkteFormat = schnittstellen.Format[
	Objekte,
	*Objekte,
]

type store struct {
	standort standort.Standort
	oaf      schnittstellen.ObjekteAkteFactory
	pool     *collections.Pool[Objekte]
	ObjekteFormat
	AkteFormat
	ObjekteInflator
	ObjekteSaver
	AkteTextSaver
}

func MakeStore(
	standort standort.Standort,
	oaf schnittstellen.ObjekteAkteFactory,
) (s *store, err error) {
	p := collections.MakePool[Objekte]()
	of := MakeFormatObjekte()
	af := MakeFormatAkte()

	s = &store{
		standort:      standort,
		oaf:           oaf,
		pool:          p,
		ObjekteFormat: of,
		AkteFormat:    af,
		ObjekteInflator: objekte.MakeObjekteInflator[
			Objekte,
			*Objekte,
			objekte.NilVerzeichnisse[Objekte],
			*objekte.NilVerzeichnisse[Objekte],
		](
			oaf,
			oaf,
			of,
			af,
			p,
		),
		ObjekteSaver: objekte.MakeObjekteSaver[
			Objekte,
			*Objekte,
		](oaf, of),
		AkteTextSaver: objekte.MakeAkteTextSaver[
			Objekte,
			*Objekte,
		](oaf, af),
	}

	return
}

func (s *store) Create(o *Objekte) (sh schnittstellen.Sha, err error) {
	if o.Akte.Skus.Len() == 0 {
		err = errors.Wrap(ErrEmpty)
		return
	}

	if s.SaveAkteText(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	if sh, err = s.SaveObjekte(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *store) readOnePath(p string) (o *Objekte, err error) {
	var sh sha.Sha

	if sh, err = sha.MakeShaFromPath(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	if o, err = s.ReadOne(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *store) ReadOne(sh sha.Sha) (o *Objekte, err error) {
	var or sha.ReadCloser

	if or, err = s.oaf.ObjekteReader(gattung.Bestandsaufnahme, sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, or)

	o = s.pool.Get()

	if _, err = s.ObjekteFormat.Parse(or, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	var ar sha.ReadCloser

	if ar, err = s.oaf.AkteReader(o.AkteSha); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ar)

	if _, err = s.AkteFormat.Parse(ar, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *store) ReadAll(f collections.WriterFunc[*Objekte]) (err error) {
	if err = files.ReadDirNamesLevel2(
		func(p string) (err error) {
			var o *Objekte

			if o, err = s.readOnePath(p); err != nil {
				err = errors.Wrap(err)
				return
			}

			shouldRepool := true

			if err = f(o); err != nil {
				if collections.IsDoNotRepool(err) {
					shouldRepool = false
					err = nil
				} else {
					err = errors.Wrap(err)
					return
				}
			}

			if shouldRepool {
				s.pool.Put(o)
			}

			return
		},
		s.standort.DirObjektenBestandsaufnahme(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
