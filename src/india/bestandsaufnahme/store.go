package bestandsaufnahme

import (
	"fmt"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/hotel/objekte"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
)

type Store interface {
	objekte_store.ObjekteSaver
	AkteTextSaver
	Create(*Akte) error
	objekte_store.LastReader[Transacted]
	objekte_store.OneReader[schnittstellen.ShaLike, *Transacted]
	objekte_store.AllReader[*Transacted]
}

type AkteTextSaver = objekte_store.AkteTextSaver[
	Akte,
	*Akte,
]

type AkteFormat = objekte.AkteFormat[
	Akte,
	*Akte,
]

type store struct {
	standort                  standort.Standort
	sv                        schnittstellen.StoreVersion
	of                        schnittstellen.ObjekteIOFactory
	af                        schnittstellen.AkteIOFactory
	pool                      schnittstellen.Pool[Akte, *Akte]
	persistentMetadateiFormat objekte_format.Format
	formatAkte
	objekte_store.ObjekteSaver
	AkteTextSaver
}

func MakeStore(
	standort standort.Standort,
	sv schnittstellen.StoreVersion,
	of schnittstellen.ObjekteIOFactory,
	af schnittstellen.AkteIOFactory,
	pmf objekte_format.Format,
) (s *store, err error) {
	p := collections.MakePool[Akte]()
	fa := formatAkte{
		af: af,
	}

	s = &store{
		standort:                  standort,
		sv:                        sv,
		of:                        of,
		af:                        af,
		pool:                      p,
		persistentMetadateiFormat: pmf,
		formatAkte:                fa,
		ObjekteSaver:              objekte_store.MakeObjekteSaver(of, pmf),
		AkteTextSaver: objekte_store.MakeAkteTextSaver[
			Akte,
			*Akte,
		](
			af,
			objekte_store.MakeAkteFormat[Akte, *Akte](
				objekte.MakeReaderAkteParseSaver[Akte, *Akte](af, fa),
				fa,
				af,
			),
		),
	}

	return
}

func (s *store) Create(o *Akte) (err error) {
	if o.Skus.Len() == 0 {
		err = errors.Wrap(ErrEmpty)
		return
	}

	var sh schnittstellen.ShaLike

	if sh, _, err = s.SaveAkteText(*o); err != nil {
		err = errors.Wrap(err)
		return
	}

	t := &Transacted{}
	t.Reset()
	t.Akte = *o
	t.SetAkteSha(sh)
	t.SetTai(kennung.NowTai())

	if err = s.SaveObjekteIncludeTai(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *store) readOnePath(p string) (o *Transacted, err error) {
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

func (s *store) ReadOne(sh schnittstellen.ShaLike) (o *Transacted, err error) {
	var or sha.ReadCloser

	if or, err = s.of.ObjekteReader(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, or)

	o = &Transacted{}
	o.Reset()

	o.SetObjekteSha(sh)

	if _, err = s.persistentMetadateiFormat.ParsePersistentMetadatei(
		or,
		o,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var ar sha.ReadCloser

	if ar, err = s.af.AkteReader(o.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ar)

	var akteSha schnittstellen.ShaLike

	if akteSha, _, err = s.formatAkte.ParseSaveAkte(ar, &o.Akte); err != nil {
		err = errors.Wrap(err)
		return
	}

	o.SetAkteSha(akteSha)

	return
}

func (s *store) ReadLast() (max Transacted, err error) {
	l := &sync.Mutex{}

	if err = s.ReadAll(
		func(b *Transacted) (err error) {
			l.Lock()
			defer l.Unlock()

			if max.Less(*b) {
				max.ResetWith(*b)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if max.GetObjekteSha().IsNull() {
		panic(fmt.Sprintf("did not find last Bestandsaufnahme: %#v", max.GetMetadatei()))
	}

	return
}

func (s *store) ReadAll(f schnittstellen.FuncIter[*Transacted]) (err error) {
	var p string

	if p, err = s.standort.DirObjektenGattung(
		s.sv,
		gattung.Bestandsaufnahme,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = files.ReadDirNamesLevel2(
		func(p string) (err error) {
			var o *Transacted

			if o, err = s.readOnePath(p); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = f(o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
		p,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
