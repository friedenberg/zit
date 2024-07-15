package object_id_index

import (
	"io"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/collections_delta"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/zettel_id_index"
)

type ObjectIdIndex[
	T ids.IdGeneric[T],
	TPtr ids.IdGenericPtr[T],
] interface {
	GetInt(int) (T, error)
	Get(*T) (*ids.IndexedLike, error)
	HasChanges() bool
	Reset() error
	GetAll() ([]ids.IdLike, error)
	Each(interfaces.FuncIter[ids.IndexedLike]) error
	EachSchwanzen(interfaces.FuncIter[*ids.IndexedLike]) error
	StoreDelta(interfaces.Delta[T]) (err error)
	StoreMany(interfaces.SetLike[T]) (err error)
	StoreOne(T) (err error)
	io.WriterTo
	io.ReaderFrom
	Flush() error
}

type EtikettIndexMutation interface {
	AddEtikettSet(to ids.TagSet, from ids.TagSet) (err error)
	Add(s ids.TagSet) (err error)
}

type EtikettIndex interface {
	EtikettIndexMutation

	EachSchwanzen(
		interfaces.FuncIter[*ids.IndexedTag],
	) error
	GetEtikett(
		*ids.Tag,
	) (*ids.IndexedLike, error)
}

type Index interface {
	interfaces.Flusher

	EtikettIndex

	zettel_id_index.HinweisIndex
}

type index struct {
	path string
	interfaces.CacheIOFactory
	didRead    bool
	hasChanges bool
	lock       *sync.RWMutex

	etikettenIndex ObjectIdIndex[ids.Tag, *ids.Tag]
	hinweisIndex   zettel_id_index.HinweisIndex
}

func MakeIndex(
	k interfaces.Config,
	s fs_home.Home,
	vf interfaces.CacheIOFactory,
) (i *index, err error) {
	i = &index{
		path:           s.FileVerzeichnisseEtiketten(),
		CacheIOFactory: vf,
		lock:           &sync.RWMutex{},
		etikettenIndex: MakeIndex2[ids.Tag](
			vf,
			s.DirVerzeichnisse("EtikettenIndexV0"),
		),
	}

	if i.hinweisIndex, err = zettel_id_index.MakeIndex(
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
	to ids.TagSet,
	from ids.TagSet,
) (err error) {
	return i.etikettenIndex.StoreDelta(
		collections_delta.MakeSetDelta[ids.Tag](from, to),
	)
}

func (i *index) GetEtikett(
	k *ids.Tag,
) (id *ids.IndexedLike, err error) {
	return i.etikettenIndex.Get(k)
}

func (i *index) Add(s ids.TagSet) (err error) {
	return i.etikettenIndex.StoreMany(s)
}

func (i *index) Each(
	f interfaces.FuncIter[ids.IndexedLike],
) (err error) {
	return i.etikettenIndex.Each(f)
}

func (i *index) EachSchwanzen(
	f interfaces.FuncIter[*ids.IndexedLike],
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

func (i *index) AddHinweis(k ids.IdLike) (err error) {
	if err = i.hinweisIndex.AddHinweis(k); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *index) CreateHinweis() (h *ids.ZettelId, err error) {
	if h, err = i.hinweisIndex.CreateHinweis(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *index) PeekHinweisen(n int) (hs []*ids.ZettelId, err error) {
	if hs, err = i.hinweisIndex.PeekHinweisen(n); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
