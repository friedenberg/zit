package kennung_index

import (
	"io"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/charlie/verzeichnisse_index"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/hinweis_index"
)

type KennungIndex[T kennung.KennungSansGattung] interface {
	Get(T) (kennung.IndexedLike[T], error)
	DidRead() bool
	HasChanges() bool
	Reset() error
	GetAll() []T
	Each(schnittstellen.FuncIter[kennung.IndexedLike[T]]) error
	EachSchwanzen(schnittstellen.FuncIter[kennung.IndexedLike[T]]) error
	StoreDelta(schnittstellen.Delta[T]) (err error)
	StoreMany(schnittstellen.Set[T]) (err error)
	StoreOne(T) (err error)
	io.WriterTo
	io.ReaderFrom
}

type EtikettIndex interface {
	Each(schnittstellen.FuncIter[kennung.IndexedLike[kennung.Etikett]]) error
	EachSchwanzen(
		schnittstellen.FuncIter[kennung.IndexedLike[kennung.Etikett]],
	) error
	AddEtikettSet(to kennung.EtikettSet, from kennung.EtikettSet) (err error)
	Add(s kennung.EtikettSet) (err error)
	GetEtikett(kennung.Etikett) (kennung.IndexedLike[kennung.Etikett], error)
}

type Index interface {
	schnittstellen.Flusher
	schnittstellen.Resetter

	EtikettIndex

	hinweis_index.HinweisIndex
}

type index struct {
	path string
	schnittstellen.VerzeichnisseFactory
	didRead    bool
	hasChanges bool
	lock       *sync.RWMutex

	etikettenIndex verzeichnisse_index.Wrapper[KennungIndex[kennung.Etikett]]
	hinweisIndex   hinweis_index.HinweisIndex
}

// TODO-P1 remove
type row struct {
	kennung.Etikett
	count int64
}

func MakeIndex(
	k schnittstellen.Konfig,
	s standort.Standort,
	vf schnittstellen.VerzeichnisseFactory,
) (i *index, err error) {
	i = &index{
		path:                 s.FileVerzeichnisseEtiketten(),
		VerzeichnisseFactory: vf,
		lock:                 &sync.RWMutex{},
		etikettenIndex: verzeichnisse_index.MakeWrapper[KennungIndex[kennung.Etikett]](
			MakeIndex2[kennung.Etikett](),
			s.DirVerzeichnisse("EtikettenIndexV0"),
		),
	}

	if i.hinweisIndex, err = hinweis_index.MakeIndex(
		k,
		s,
		vf,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *index) Flush() (err error) {
	if err = i.hinweisIndex.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = i.etikettenIndex.Flush(i.VerzeichnisseFactory); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *index) AddEtikettSet(
	to kennung.EtikettSet,
	from kennung.EtikettSet,
) (err error) {
	var ei KennungIndex[kennung.Etikett]

	if ei, err = i.etikettenIndex.Get(i.VerzeichnisseFactory); err != nil {
		err = errors.Wrap(err)
		return
	}

	return ei.StoreDelta(collections.MakeSetDelta(from, to))
}

func (i *index) GetEtikett(
	k kennung.Etikett,
) (id kennung.IndexedLike[kennung.Etikett], err error) {
	var ei KennungIndex[kennung.Etikett]

	if ei, err = i.etikettenIndex.Get(i.VerzeichnisseFactory); err != nil {
		err = errors.Wrap(err)
		return
	}

	return ei.Get(k)
}

func (i *index) Add(s kennung.EtikettSet) (err error) {
	var ei KennungIndex[kennung.Etikett]

	if ei, err = i.etikettenIndex.Get(i.VerzeichnisseFactory); err != nil {
		err = errors.Wrap(err)
		return
	}

	return ei.StoreMany(s)
}

func (i *index) Each(
	f schnittstellen.FuncIter[kennung.IndexedLike[kennung.Etikett]],
) (err error) {
	var ei KennungIndex[kennung.Etikett]

	if ei, err = i.etikettenIndex.Get(i.VerzeichnisseFactory); err != nil {
		err = errors.Wrap(err)
		return
	}

	return ei.Each(f)
}

func (i *index) EachSchwanzen(
	f schnittstellen.FuncIter[kennung.IndexedLike[kennung.Etikett]],
) (err error) {
	var ei KennungIndex[kennung.Etikett]

	if ei, err = i.etikettenIndex.Get(i.VerzeichnisseFactory); err != nil {
		err = errors.Wrap(err)
		return
	}

	return ei.EachSchwanzen(f)
}

func (i *index) Reset() (err error) {
	if err = i.hinweisIndex.Reset(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var ei KennungIndex[kennung.Etikett]

	if ei, err = i.etikettenIndex.Get(i.VerzeichnisseFactory); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = ei.Reset(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *index) AddHinweis(h kennung.Hinweis) (err error) {
	if err = i.hinweisIndex.AddHinweis(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *index) CreateHinweis() (h kennung.Hinweis, err error) {
	if h, err = i.hinweisIndex.CreateHinweis(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *index) PeekHinweisen(n int) (hs []kennung.Hinweis, err error) {
	if hs, err = i.hinweisIndex.PeekHinweisen(n); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
