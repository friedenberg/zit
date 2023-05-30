package kennung_index

import (
	"bufio"
	"encoding/gob"
	"io"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type index2[
	T kennung.KennungLike[T],
	TPtr kennung.KennungLikePtr[T],
] struct {
	didRead    bool
	hasChanges bool
	lock       *sync.Mutex
	Kennungen  map[string]indexed[T, TPtr]
}

func MakeIndex2[
	T kennung.KennungLike[T],
	TPtr kennung.KennungLikePtr[T],
]() (i *index2[T, TPtr]) {
	i = &index2[T, TPtr]{
		lock:      &sync.Mutex{},
		Kennungen: make(map[string]indexed[T, TPtr]),
	}

	return
}

func (i *index2[T, TPtr]) DidRead() bool {
	i.lock.Lock()
	defer i.lock.Unlock()

	return i.didRead
}

func (i *index2[T, TPtr]) HasChanges() bool {
	i.lock.Lock()
	defer i.lock.Unlock()

	return i.hasChanges
}

func (i *index2[T, TPtr]) Reset() error {
	i.lock.Lock()
	defer i.lock.Unlock()

	i.Kennungen = make(map[string]indexed[T, TPtr])

	return nil
}

func (i index2[T, TPtr]) WriteTo(w1 io.Writer) (n int64, err error) {
	w := bufio.NewWriter(w1)
	defer errors.DeferredFlusher(&err, w)

	i.lock.Lock()
	defer i.lock.Unlock()

	enc := gob.NewEncoder(w)

	if err = enc.Encode(i.Kennungen); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *index2[T, TPtr]) ReadFrom(r1 io.Reader) (n int64, err error) {
	r := bufio.NewReader(r1)

	dec := gob.NewDecoder(r)

	i.lock.Lock()
	defer i.lock.Unlock()

	if err = dec.Decode(&i.Kennungen); err != nil {
		if errors.IsEOF(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	i.didRead = true

	return
}

func (i index2[T, TPtr])	GetAll() (out []T) {
  out = make([]T, 0, len(i.Kennungen))

  for _, ki := range i.Kennungen {
    out = append(out, ki.GetKennung())
  }

  return
}

func (i index2[T, TPtr]) Get(k T) (id Indexed2[T], ok bool) {
	i.lock.Lock()
	defer i.lock.Unlock()

	id, ok = i.Kennungen[k.String()]
	return
}

func (i *index2[T, TPtr]) StoreMany(ks schnittstellen.Set[T]) (err error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	return ks.Each(i.storeOne)
}

func (i *index2[T, TPtr]) StoreOne(k T) (err error) {
	if kennung.IsEmpty(k) {
		return
	}

	i.lock.Lock()
	defer i.lock.Unlock()

	return i.storeOne(k)
}

func (i *index2[T, TPtr]) storeOne(k T) (err error) {
	if _, ok := i.Kennungen[k.String()]; ok {
		return
	}

	i.hasChanges = true

	id := indexed[T, TPtr]{}
	id.ResetWithKennung(k)
	i.Kennungen[k.String()] = id
	return
}
