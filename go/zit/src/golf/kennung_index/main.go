package kennung_index

import (
	"io"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/collections_delta"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/echo/standort"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/hinweis_index"
)

type KennungIndex[
	T kennung.KennungLike[T],
	TPtr kennung.KennungLikePtr[T],
] interface {
	GetInt(int) (T, error)
	Get(*T) (*kennung.IndexedLike, error)
	HasChanges() bool
	Reset() error
	GetAll() ([]kennung.Kennung, error)
	Each(interfaces.FuncIter[kennung.IndexedLike]) error
	EachSchwanzen(interfaces.FuncIter[*kennung.IndexedLike]) error
	StoreDelta(interfaces.Delta[T]) (err error)
	StoreMany(interfaces.SetLike[T]) (err error)
	StoreOne(T) (err error)
	io.WriterTo
	io.ReaderFrom
	Flush() error
}

type EtikettIndexMutation interface {
	AddEtikettSet(to kennung.EtikettSet, from kennung.EtikettSet) (err error)
	Add(s kennung.EtikettSet) (err error)
}

type EtikettIndex interface {
	EtikettIndexMutation

	EachSchwanzen(
		interfaces.FuncIter[*kennung.IndexedEtikett],
	) error
	GetEtikett(
		*kennung.Etikett,
	) (*kennung.IndexedLike, error)
}

type Index interface {
	interfaces.Flusher

	EtikettIndex

	hinweis_index.HinweisIndex
}

type index struct {
	path string
	interfaces.VerzeichnisseFactory
	didRead    bool
	hasChanges bool
	lock       *sync.RWMutex

	etikettenIndex KennungIndex[kennung.Etikett, *kennung.Etikett]
	hinweisIndex   hinweis_index.HinweisIndex
}

func MakeIndex(
	k interfaces.Konfig,
	s standort.Standort,
	vf interfaces.VerzeichnisseFactory,
) (i *index, err error) {
	i = &index{
		path:                 s.FileVerzeichnisseEtiketten(),
		VerzeichnisseFactory: vf,
		lock:                 &sync.RWMutex{},
		etikettenIndex: MakeIndex2[kennung.Etikett](
			vf,
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

	if err = i.etikettenIndex.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *index) AddEtikettSet(
	to kennung.EtikettSet,
	from kennung.EtikettSet,
) (err error) {
	return i.etikettenIndex.StoreDelta(
		collections_delta.MakeSetDelta[kennung.Etikett](from, to),
	)
}

func (i *index) GetEtikett(
	k *kennung.Etikett,
) (id *kennung.IndexedLike, err error) {
	return i.etikettenIndex.Get(k)
}

func (i *index) Add(s kennung.EtikettSet) (err error) {
	return i.etikettenIndex.StoreMany(s)
}

func (i *index) Each(
	f interfaces.FuncIter[kennung.IndexedLike],
) (err error) {
	return i.etikettenIndex.Each(f)
}

func (i *index) EachSchwanzen(
	f interfaces.FuncIter[*kennung.IndexedLike],
) (err error) {
	return i.etikettenIndex.EachSchwanzen(f)
}

func (i *index) Reset() (err error) {
	if err = i.hinweisIndex.Reset(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = i.etikettenIndex.Reset(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *index) AddHinweis(k kennung.Kennung) (err error) {
	if err = i.hinweisIndex.AddHinweis(k); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *index) CreateHinweis() (h *kennung.Hinweis, err error) {
	if h, err = i.hinweisIndex.CreateHinweis(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *index) PeekHinweisen(n int) (hs []*kennung.Hinweis, err error) {
	if hs, err = i.hinweisIndex.PeekHinweisen(n); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
