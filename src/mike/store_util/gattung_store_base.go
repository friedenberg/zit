package store_util

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/kilo/objekte_store"
)

type GattungStoreLike interface {
	Reindexer
}

type CommonStoreBase[
	O schnittstellen.Akte[O],
	OPtr schnittstellen.AktePtr[O],
] struct {
	schnittstellen.GattungGetter
	StoreUtil

	delegate CommonStoreDelegate

	objekte_store.TransactedReader
	objekte_store.LogWriter
	persistentMetadateiFormat objekte_format.Format

	// objekte_store.AkteTextSaver[O, OPtr]
	// objekte_store.StoredParseSaver[O, OPtr]
	akteFormat objekte.AkteFormat[O, OPtr]
}

func MakeCommonStoreBase[
	O schnittstellen.Akte[O],
	OPtr schnittstellen.AktePtr[O],
](
	gg schnittstellen.GattungGetter,
	delegate CommonStoreDelegate,
	sa StoreUtil,
	tr objekte_store.TransactedReader,
	pmf objekte_format.Format,
	akteFormat objekte.AkteFormat[O, OPtr],
) (s *CommonStoreBase[O, OPtr], err error) {
	if delegate == nil {
		panic("delegate was nil")
	}

	s = &CommonStoreBase[O, OPtr]{
		GattungGetter:             gg,
		delegate:                  delegate,
		StoreUtil:                 sa,
		akteFormat:                akteFormat,
		TransactedReader:          tr,
		persistentMetadateiFormat: pmf,
	}

	return
}

func (s *CommonStoreBase[O, OPtr]) SetLogWriter(
	lw objekte_store.LogWriter,
) {
	s.LogWriter = lw
}

func (s *CommonStoreBase[O, OPtr]) Query(
	m matcher.MatcherSigil,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	return objekte_store.QueryMethodForMatcher(s, m, f)
}

func (s *CommonStoreBase[O, OPtr]) ReindexOne(
	t *sku.Transacted,
) (o matcher.Matchable, err error) {
	o = t

	if t.IsNew() {
		s.LogWriter.New(t)
		if err = s.delegate.AddOne(t); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		s.LogWriter.Updated(t)
		if err = s.delegate.UpdateOne(t); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
