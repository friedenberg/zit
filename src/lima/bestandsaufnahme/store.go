package bestandsaufnahme

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/log"
	"github.com/friedenberg/zit/src/bravo/pool"
	"github.com/friedenberg/zit/src/charlie/catgut"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/sku_fmt"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/kilo/objekte_store"
)

type Store interface {
	GetStore() Store
	objekte_store.AkteTextSaver[Akte, *Akte]
	Create(*Akte) (*sku.Transacted, error)
	objekte_store.LastReader
	ReadOne(schnittstellen.Stringer) (*sku.Transacted, error)
	ReadOneSku(besty, sk *sha.Sha) (*sku.Transacted, error)
	objekte_store.AllReader
	ReadAllSkus(func(besty, sk *sku.Transacted) error) error
	schnittstellen.AkteGetter[*Akte]
}

type AkteFormat = objekte.AkteFormat[
	Akte,
	*Akte,
]

type akteFormat interface {
	FormatParsedAkte(io.Writer, *Akte) (n int64, err error)
	objekte.AkteParser[Akte, *Akte]
}

type store struct {
	standort                  standort.Standort
	ls                        schnittstellen.LockSmith
	sv                        schnittstellen.StoreVersion
	of                        schnittstellen.ObjekteIOFactory
	af                        schnittstellen.AkteIOFactory
	clock                     kennung.Clock
	pool                      schnittstellen.Pool[Akte, *Akte]
	persistentMetadateiFormat objekte_format.Format
	options                   objekte_format.Options
	formatAkte                akteFormat
	objekte_store.AkteTextSaver[Akte, *Akte]
}

func MakeStore(
	standort standort.Standort,
	ls schnittstellen.LockSmith,
	sv schnittstellen.StoreVersion,
	of schnittstellen.ObjekteIOFactory,
	af schnittstellen.AkteIOFactory,
	pmf objekte_format.Format,
	clock kennung.Clock,
) (s *store, err error) {
	p := pool.MakePool[Akte](nil, func(a *Akte) { Resetter.Reset(a) })

	op := objekte_format.Options{IncludeTai: true}
	fa := MakeAkteFormat(sv, op)

	s = &store{
		standort:                  standort,
		ls:                        ls,
		sv:                        sv,
		of:                        of,
		af:                        af,
		pool:                      p,
		clock:                     clock,
		persistentMetadateiFormat: pmf,
		options:                   op,
		formatAkte:                fa,
		AkteTextSaver: objekte_store.MakeAkteStore[
			Akte,
			*Akte,
		](
			standort,
			objekte_store.MakeAkteFormat[Akte, *Akte](
				objekte.MakeReaderAkteParseSaver[Akte, *Akte](af, fa),
				fa,
				af,
			),
			Resetter.Reset,
		),
	}

	return
}

func (s *store) GetStore() Store {
	return s
}

func (s *store) Create(o *Akte) (t *sku.Transacted, err error) {
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

	if sh, _, err = s.SaveAkteText(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	t = sku.GetTransactedPool().Get()
	defer sku.GetTransactedPool().Put(t)

	sku.TransactedResetter.Reset(t)
	t.SetAkteSha(sh)
	tai := s.clock.GetTai()

	if err = t.Kennung.SetWithKennung(tai); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.SetTai(tai)

	var w sha.WriteCloser

	if w, err = s.of.ObjekteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, w)

	if _, err = s.persistentMetadateiFormat.FormatPersistentMetadatei(
		w,
		t,
		objekte_format.Options{IncludeTai: true},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = sha.Make(w.GetShaLike())

	log.Log().Printf(
		"saving Bestandsaufnahme with tai: %s -> %s",
		t.GetKennung().GetGattung(),
		sh,
	)

	t.SetObjekteSha(sh)

	return
}

func (s *store) readOnePath(p string) (o *sku.Transacted, err error) {
	var sh *sha.Sha

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

func (s *store) ReadOneSku(besty, sh *sha.Sha) (o *sku.Transacted, err error) {
	var bestySku *sku.Transacted

	if bestySku, err = s.ReadOne(besty); err != nil {
		err = errors.Wrap(err)
		return
	}

	var ar schnittstellen.ShaReadCloser

	if ar, err = s.af.AkteReader(&bestySku.Metadatei.Akte); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ar)

	dec := sku_fmt.MakeFormatBestandsaufnahmeScanner(
		ar,
		s.persistentMetadateiFormat,
		s.options,
	)

	for dec.Scan() {
		o = dec.GetTransacted()

		if !o.Metadatei.Sha.Equals(sh) {
			continue
		}

		return
	}

	o = nil

	if err = dec.Error(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *store) ReadOne(
	k schnittstellen.Stringer,
) (o *sku.Transacted, err error) {
	var sh sha.Sha

	if err = sh.Set(k.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	var or sha.ReadCloser

	if or, err = s.of.ObjekteReader(&sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, or)

	var sb strings.Builder
	mr := io.TeeReader(or, &sb)

	var n int64
	n, o, err = s.readOneFromReader(mr)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	switch s.sv.GetInt() {
	case 0, 1, 2:
		panic("unsupported version")

	default:
		op := s.options
		err = sku.CalculateAndConfirmSha(
			o,
			s.persistentMetadateiFormat,
			op,
			&sh,
		)

		if err != nil {
			op.IncludeTai = false

			err = sku.CalculateAndConfirmSha(
				o,
				s.persistentMetadateiFormat,
				op,
				&sh,
			)

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		if err != nil {
			err = errors.Wrapf(err, "Sku: %q, Read: %q, N: %d", o, sb.String(), n)
			return
		}
	}

	return
}

func (s *store) readOneFromReader(
	r io.Reader,
) (n int64, o *sku.Transacted, err error) {
	o = sku.GetTransactedPool().Get()

	if n, err = s.persistentMetadateiFormat.ParsePersistentMetadatei(
		catgut.MakeRingBuffer(r, 0),
		o,
		s.options,
	); err != nil {
		if errors.IsEOF(err) {
			err = nil
		} else {
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

	a = MakeAkte()

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

func (s *store) ReadLast() (max *sku.Transacted, err error) {
	l := &sync.Mutex{}

	var maxSku sku.Transacted

	if err = s.ReadAll(
		func(b *sku.Transacted) (err error) {
			l.Lock()
			defer l.Unlock()

			if sku.TransactedLessor.LessPtr(&maxSku, b) {
				sku.TransactedResetter.ResetWith(&maxSku, b)
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

func (s *store) ReadAll(
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
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
			var o *sku.Transacted

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
	f func(besty, sk *sku.Transacted) error,
) (err error) {
	if err = s.ReadAll(
		func(t *sku.Transacted) (err error) {
			var a *Akte

			if a, err = s.GetAkte(t.GetAkteSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = a.Skus.EachPtr(
				func(sk *sku.Transacted) (err error) {
					return f(t, sk)
				},
			); err != nil {
				err = errors.Wrapf(
					err,
					"Bestandsaufnahme: %s",
					t.GetKennung(),
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

func (s *store) ReadAllSkus2(
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	if err = s.ReadAll(
		func(t *sku.Transacted) (err error) {
			var r io.ReadCloser

			if r, err = s.af.AkteReader(t.GetAkteSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer errors.DeferredCloser(&err, r)

			dec := sku_fmt.MakeFormatBestandsaufnahmeScanner(
				r,
				s.persistentMetadateiFormat,
				s.options,
			)

			for dec.Scan() {
				sk := dec.GetTransacted()

				if err = f(sk); err != nil {
					err = errors.Wrapf(err, "Sku: %s", sk)
					return
				}
			}

			if err = dec.Error(); err != nil {
				err = errors.Wrap(err)
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
