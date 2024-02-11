package bestandsaufnahme

import (
	"fmt"
	"io"
	"sort"
	"sync"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/files"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/bravo/log"
	"code.linenisgreat.com/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/src/charlie/catgut"
	"code.linenisgreat.com/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/charlie/sha"
	"code.linenisgreat.com/zit/src/delta/standort"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/golf/ennui_shas"
	"code.linenisgreat.com/zit/src/golf/objekte_format"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/src/juliett/objekte"
	"code.linenisgreat.com/zit/src/kilo/objekte_store"
)

type Store interface {
	errors.Flusher
	GetStore() Store
	ReadOneEnnui(*sha.Sha) (*sku.Transacted, error)
	ReadOneKennung(kennung.Kennung) (*sku.Transacted, error)
	ReadOneKennungSha(kennung.Kennung) (*sha.Sha, error)
	Create(*Akte) (*sku.Transacted, error)
	objekte_store.LastReader
	WriteOneObjekteMetadatei(o *sku.Transacted) (err error)
	ReadOne(schnittstellen.Stringer) (*sku.Transacted, error)
	ReadOneSku(besty, sk *sha.Sha) (*sku.Transacted, error)
	objekte_store.AllReader
	ReadAllSkus(func(besty, sk *sku.Transacted) error) error
	schnittstellen.AkteGetter[*Akte]

	StreamAkte(
		schnittstellen.ShaLike,
		schnittstellen.FuncIter[*sku.Transacted],
	) error
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
	ennuiKennung              ennui_shas.Ennui
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

	op := objekte_format.Options{Tai: true}
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
	}

	if s.ennuiKennung, err = ennui_shas.MakeNoDuplicates(
		standort,
		standort.DirVerzeichnisseVerweise(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *store) GetStore() Store {
	return s
}

func (s *store) Flush() (err error) {
	wg := iter.MakeErrorWaitGroupParallel()
	wg.Do(s.ennuiKennung.Flush)
	return wg.GetError()
}

func (s *store) ReadOneEnnui(sh *sha.Sha) (sk *sku.Transacted, err error) {
	var r sha.ReadCloser

	if r, err = s.standort.AkteReaderFrom(
		sh,
		s.standort.DirVerzeichnisseMetadateiKennungMutter(),
	); err != nil {
		if errors.IsNotExist(err) {
			err = collections.MakeErrNotFound(sh)
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	defer errors.DeferredCloser(&err, r)

	rb := catgut.MakeRingBuffer(r, 0)

	sk = sku.GetTransactedPool().Get()

	var n int64

	n, err = s.persistentMetadateiFormat.ParsePersistentMetadatei(
		rb,
		sk,
		s.options,
	)

	if err == io.EOF && n > 0 {
		err = nil
	} else if err != io.EOF && err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = sk.CalculateObjekteShas(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *store) ReadOneKennungSha(k kennung.Kennung) (sh *sha.Sha, err error) {
	left := sha.FromString(k.String())
	defer sha.GetPool().Put(left)

	if sh, err = s.ennuiKennung.ReadOne(left); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *store) ReadOneKennung(k kennung.Kennung) (sk *sku.Transacted, err error) {
	sh, err := s.ReadOneKennungSha(k)
	defer sha.GetPool().Put(sh)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if sk, err = s.ReadOneEnnui(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !sh.Equals(sk.Metadatei.Sha()) {
		err = errors.Errorf(
			"expected sha %q but got %q",
			sh,
			sk.Metadatei.Sha(),
		)

		return
	}

	return
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

	var sh *sha.Sha

	if sh, err = s.writeAkte(o); err != nil {
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
		objekte_format.Options{Tai: true},
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

func (s *store) writeAkte(o *Akte) (sh *sha.Sha, err error) {
	var sw sha.WriteCloser

	if sw, err = s.standort.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, sw)

	fo := sku_fmt.MakeFormatBestandsaufnahmePrinter(
		sw,
		s.persistentMetadateiFormat,
		s.options,
	)

	defer o.Skus.Restore()

	wg := iter.MakeErrorWaitGroupParallel()
	tickets := make(chan struct{}, 128)

	for {
		sk, ok := o.Skus.PopAndSave()

		if !ok {
			break
		}

		if sk.Metadatei.Sha().IsNull() {
			err = errors.Errorf("empty sha: %s", sk)
			return
		}

		_, err = fo.Print(sk)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		wg.Do(
			func() (err error) {
				tickets <- struct{}{}
				defer func() {
					<-tickets
				}()

				if err = s.WriteOneObjekteMetadatei(sk); err != nil {
					err = errors.Wrap(err)
					return
				}

				return
			},
		)
	}

	sh = sha.Make(sw.GetShaLike())

	if err = wg.GetError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *store) makeWriteMetadateiFunc(
	dir string,
	fo objekte_format.FormatGeneric,
	o *sku.Transacted,
	expected *sha.Sha,
) schnittstellen.FuncError {
	return func() (err error) {
		var sw sha.WriteCloser

		if sw, err = s.standort.AkteWriterToLight(dir); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, sw)

		if _, err = fo.WriteMetadateiTo(sw, o); err != nil {
			err = errors.Wrap(err)
			return
		}

		actual := sw.GetShaLike()

		if !expected.EqualsSha(actual) {
			err = errors.Errorf(
				"expected %q but got %q",
				expected,
				actual,
			)

			return
		}

		return
	}
}

func (s *store) WriteOneObjekteMetadatei(o *sku.Transacted) (err error) {
	if o.Metadatei.Sha().IsNull() {
		err = errors.Errorf("null sha")
		return
	}

	wg := iter.MakeErrorWaitGroupParallel()

	wg.Do(
		func() (err error) {
			k := o.GetKennung()
			sh := sha.FromString(k.String())
			defer sha.GetPool().Put(sh)

			if err = s.ennuiKennung.AddSha(sh, o.Metadatei.Sha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	)

	wg.Do(
		s.makeWriteMetadateiFunc(
			s.standort.DirVerzeichnisseMetadateiKennungMutter(),
			objekte_format.Formats.MetadateiKennungMutter(),
			o,
			o.Metadatei.Sha(),
		),
	)

	wg.Do(
		s.makeWriteMetadateiFunc(
			s.standort.DirVerzeichnisseMetadatei(),
			objekte_format.Formats.Metadatei(),
			o,
			&o.Metadatei.SelbstMetadatei,
		),
	)

	return wg.GetError()
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

	if err = o.CalculateObjekteShas(); err != nil {
		if errors.Is(err, objekte_format.ErrEmptyTai) {
			var t kennung.Tai
			err1 := t.Set(o.Kennung.String())

			if err1 != nil {
				err = errors.Wrapf(err, "%#v", o)
				return
			}

			o.SetTai(t)

			if err = o.CalculateObjekteShas(); err != nil {
				err = errors.Wrapf(err, "%#v", o)
				return
			}
		} else {
			err = errors.Wrapf(err, "%#v", o)
		}

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

		if !o.Metadatei.Sha().Equals(sh) {
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

	_, o, err = s.readOneFromReader(or)

	if err != nil {
		err = errors.Wrap(err)
		return
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

func (s *store) populateAkte(akteSha schnittstellen.ShaLike, a *Akte) (err error) {
	var ar schnittstellen.ShaReadCloser

	if ar, err = s.af.AkteReader(akteSha); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ar)

	sw := sha.MakeWriter(nil)

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

func (s *store) StreamAkte(
	akteSha schnittstellen.ShaLike,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	var ar schnittstellen.ShaReadCloser

	if ar, err = s.af.AkteReader(akteSha); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ar)

	dec := sku_fmt.MakeFormatBestandsaufnahmeScanner(
		ar,
		objekte_format.FormatForVersion(s.sv),
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
}

func (s *store) GetAkte(akteSha schnittstellen.ShaLike) (a *Akte, err error) {
	a = MakeAkte()
	err = s.populateAkte(akteSha, a)
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

	max = &maxSku

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

func (s *store) ReadAllSorted(
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	var skus []*sku.Transacted

	if err = s.ReadAll(
		func(o *sku.Transacted) (err error) {
			skus = append(skus, o)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	sort.Slice(skus, func(i, j int) bool { return skus[i].Less(skus[j]) })

	for _, o := range skus {
		if err = f(o); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *store) ReadAllSkus(
	f func(besty, sk *sku.Transacted) error,
) (err error) {
	if err = s.ReadAll(
		func(t *sku.Transacted) (err error) {
			if err = s.StreamAkte(
				t.GetAkteSha(),
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
