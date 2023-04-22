package bestandsaufnahme

import (
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/golf/persisted_metadatei_format"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
)

type Store interface {
	ObjekteInflator
	ObjekteSaver
	AkteTextSaver
	Create(*Objekte) (schnittstellen.Sha, error)
	objekte_store.LastReader[*Objekte]
	objekte_store.OneReader[schnittstellen.Sha, *Objekte]
	objekte_store.AllReader[*Objekte]
}

type ObjekteSaver = objekte_store.ObjekteSaver[
	Objekte,
	*Objekte,
]

type AkteTextSaver = objekte_store.AkteTextSaver[
	Objekte,
	*Objekte,
]

type ObjekteInflator = objekte_store.ObjekteInflator[
	Objekte,
	*Objekte,
	objekte.NilVerzeichnisse[Objekte],
	*objekte.NilVerzeichnisse[Objekte],
]

type AkteFormat = objekte_store.AkteFormat[
	Objekte,
	*Objekte,
]

type store struct {
	standort                  standort.Standort
	of                        schnittstellen.ObjekteIOFactory
	af                        schnittstellen.AkteIOFactory
	pool                      schnittstellen.Pool[Objekte, *Objekte]
	persistentMetadateiFormat persisted_metadatei_format.Format
	formatAkte
	ObjekteInflator
	ObjekteSaver
	AkteTextSaver
}

func MakeStore(
	standort standort.Standort,
	of schnittstellen.ObjekteIOFactory,
	af schnittstellen.AkteIOFactory,
	pmf persisted_metadatei_format.Format,
) (s *store, err error) {
	p := collections.MakePool[Objekte]()
	fa := formatAkte{}

	s = &store{
		standort:                  standort,
		of:                        of,
		af:                        af,
		pool:                      p,
		persistentMetadateiFormat: pmf,
		formatAkte:                fa,
		ObjekteInflator: objekte_store.MakeObjekteInflator[
			Objekte,
			*Objekte,
			objekte.NilVerzeichnisse[Objekte],
			*objekte.NilVerzeichnisse[Objekte],
		](
			of,
			af,
			pmf,
			fa,
			p,
		),
		ObjekteSaver: objekte_store.MakeObjekteSaver[
			Objekte,
			*Objekte,
		](of, pmf),
		AkteTextSaver: objekte_store.MakeAkteTextSaver[
			Objekte,
			*Objekte,
		](
			af,
			objekte_store.MakeAkteFormat[Objekte, *Objekte](
				objekte.MakeReaderAkteParseSaver[Objekte, *Objekte](af, fa),
				fa,
				af,
			),
		),
	}

	return
}

func (s *store) Create(o *Objekte) (sh schnittstellen.Sha, err error) {
	if o.Akte.Skus.Len() == 0 {
		err = errors.Wrap(ErrEmpty)
		return
	}

	if s.SaveAkteText(*o); err != nil {
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

func (s *store) ReadOne(sh schnittstellen.Sha) (o *Objekte, err error) {
	var or sha.ReadCloser

	if or, err = s.of.ObjekteReader(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, or)

	o = s.pool.Get()

	if _, err = s.persistentMetadateiFormat.Parse(or, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	var ar sha.ReadCloser

	if ar, err = s.af.AkteReader(o.AkteSha); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ar)

	if _, err = s.formatAkte.ParseAkte(ar, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *store) ReadLast() (max *Objekte, err error) {
	l := &sync.Mutex{}

	if err = s.ReadAll(
		func(b *Objekte) (err error) {
			l.Lock()
			defer l.Unlock()

			if max == nil || max.Less(*b) {
				if max != nil {
					errors.TodoP3("repool max")
				}

				max = b
				err = collections.ErrDoNotRepool
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *store) ReadAll(f schnittstellen.FuncIter[*Objekte]) (err error) {
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
