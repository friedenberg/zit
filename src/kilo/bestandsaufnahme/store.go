package bestandsaufnahme

import (
	"fmt"
	"io"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/pool"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/lima/objekte_store"
)

type Store interface {
	objekte_store.ObjekteSaver
	AkteTextSaver
	Create(*Akte) error
	objekte_store.LastReader[*Transacted]
	objekte_store.OneReader[schnittstellen.ShaLike, *Transacted]
	objekte_store.AllReader[*Transacted]
	ReadAllSkus(schnittstellen.FuncIter[sku.SkuLikePtr]) error
	schnittstellen.AkteGetter[*Akte]
}

type AkteTextSaver = objekte_store.AkteTextSaver[
	Akte,
	*Akte,
]

type AkteFormat = objekte.AkteFormat[
	Akte,
	*Akte,
]

type akteFormat interface {
	FormatParsedAkte(io.Writer, Akte) (n int64, err error)
	objekte.AkteParser[*Akte]
}

type store struct {
	standort                  standort.Standort
	ls                        schnittstellen.LockSmith
	sv                        schnittstellen.StoreVersion
	of                        schnittstellen.ObjekteIOFactory
	af                        schnittstellen.AkteIOFactory
	pool                      schnittstellen.Pool[Akte, *Akte]
	persistentMetadateiFormat objekte_format.Format
	formatAkte                akteFormat
	objekte_store.ObjekteSaver
	AkteTextSaver
}

func MakeStore(
	standort standort.Standort,
	ls schnittstellen.LockSmith,
	sv schnittstellen.StoreVersion,
	orfg schnittstellen.ObjekteReaderFactoryGetter,
	of schnittstellen.ObjekteIOFactory,
	af schnittstellen.AkteIOFactory,
	pmf objekte_format.Format,
) (s *store, err error) {
	p := pool.MakePoolWithReset[Akte]()

	var fa akteFormat

	switch sv.GetInt() {
	case 3:
		fa = MakeAkteFormat(sv, af)

	default:
		fa = formatAkte{
			orfg:                      orfg,
			persistentMetadateiFormat: pmf,
			af:                        af,
		}
	}

	s = &store{
		standort:                  standort,
		ls:                        ls,
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
	if !s.ls.IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "create bestandsaufnahme",
		}

		return
	}

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
	t.SetAkteSha(sh)
	// TODO-P2 switch to clock
	tai := kennung.NowTai()

	t.Kennung = tai
	t.SetTai(tai)

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

func (s *store) ReadOne(
	sh schnittstellen.ShaLike,
) (o *Transacted, err error) {
	var or sha.ReadCloser

	if or, err = s.of.ObjekteReader(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, or)

	o = &Transacted{}
	o.Reset()

	if _, err = s.persistentMetadateiFormat.ParsePersistentMetadatei(
		or,
		o,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	switch s.sv.GetInt() {
	case 0, 1, 2:
		o.SetObjekteSha(sh)

	default:
		if err = sku.CalculateAndConfirmSha(o, s.persistentMetadateiFormat, sh); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *store) GetAkte(akteSha schnittstellen.ShaLike) (a *Akte, err error) {
	var ar schnittstellen.ShaReadCloser

	if ar, err = s.af.AkteReader(akteSha); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ar)

	sw := sha.MakeWriter(nil)

	a = &Akte{}
	a.Reset()

	if _, err = s.formatAkte.ParseAkte(io.TeeReader(ar, sw), a); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh := sw.GetShaLike()

	if !sh.EqualsSha(akteSha) {
		err = errors.Errorf(
			"objekte had akte sha %s while akte reader had %s",
			akteSha,
			sh,
		)
		return
	}

	return
}

func (s *store) ReadLast() (max *Transacted, err error) {
	l := &sync.Mutex{}

	var maxSku Sku

	if err = s.ReadAll(
		func(b *Transacted) (err error) {
			l.Lock()
			defer l.Unlock()

			if maxSku.Less(*b) {
				maxSku.ResetWith(*b)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if max, err = s.ReadOne(maxSku.GetObjekteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if max.GetObjekteSha().IsNull() {
		panic(
			fmt.Sprintf(
				"did not find last Bestandsaufnahme: %#v",
				max.GetMetadatei(),
			),
		)
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

// TODO-P3 support streaming instead of reading into heaps
func (s *store) ReadAllSkus(
	f schnittstellen.FuncIter[sku.SkuLikePtr],
) (err error) {
	if err = s.ReadAll(
		func(t *Transacted) (err error) {
			var a *Akte

			if a, err = s.GetAkte(t.GetAkteSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = sku.HeapEachPtr(a.Skus, f); err != nil {
				err = errors.Wrapf(
					err,
					"Bestandsaufnahme: %s",
					t.GetKennungLike(),
				)

				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}