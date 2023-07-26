package kennung_index

import (
	"bufio"
	"encoding/gob"
	"io"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/log"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type index2[
	T kennung.KennungLike[T],
	TPtr kennung.KennungLikePtr[T],
] struct {
	didRead         bool
	hasChanges      bool
	lock            *sync.RWMutex
	IntsToKennungen map[int]T
	Kennungen       map[string]kennung.Indexed[T, TPtr]
}

func MakeIndex2[
	T kennung.KennungLike[T],
	TPtr kennung.KennungLikePtr[T],
]() (i *index2[T, TPtr]) {
	i = &index2[T, TPtr]{
		lock:            &sync.RWMutex{},
		IntsToKennungen: make(map[int]T),
		Kennungen:       make(map[string]kennung.Indexed[T, TPtr]),
	}

	return
}

func (i *index2[T, TPtr]) DidRead() bool {
	i.lock.RLock()
	defer i.lock.RUnlock()

	return i.didRead
}

func (i *index2[T, TPtr]) HasChanges() bool {
	i.lock.RLock()
	defer i.lock.RUnlock()

	return i.hasChanges
}

func (i *index2[T, TPtr]) Reset() error {
	i.lock.Lock()
	defer i.lock.Unlock()

	i.Kennungen = make(map[string]kennung.Indexed[T, TPtr])

	return nil
}

func (i index2[T, TPtr]) WriteTo(w1 io.Writer) (n int64, err error) {
	w := bufio.NewWriter(w1)
	defer errors.DeferredFlusher(&err, w)

	i.lock.RLock()
	defer i.lock.RUnlock()

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

func (i index2[T, TPtr]) Each(
	f schnittstellen.FuncIter[kennung.IndexedLike[T]],
) (err error) {
	for _, id := range i.Kennungen {
		if err = f(id); err != nil {
			if iter.IsStopIteration(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (i index2[T, TPtr]) EachSchwanzen(
	f schnittstellen.FuncIter[kennung.IndexedLike[T]],
) (err error) {
	for _, id := range i.Kennungen {
		if id.GetSchwanzenCount() == 0 {
			continue
		}

		if err = f(id); err != nil {
			if iter.IsStopIteration(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (i index2[T, TPtr]) GetAll() (out []T) {
	out = make([]T, 0, len(i.Kennungen))

	for _, ki := range i.Kennungen {
		out = append(out, ki.GetKennung())
	}

	return
}

func (i index2[T, TPtr]) GetInt(in int) (id T, err error) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	ok := false
	id, ok = i.IntsToKennungen[in]

	if !ok {
		err = errors.Wrap(collections.ErrNotFound{})
		return
	}

	return
}

func (i index2[T, TPtr]) Get(k T) (id kennung.IndexedLike[T], err error) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	ok := false
	id, ok = i.Kennungen[k.String()]

	if !ok {
		err = errors.Wrap(collections.ErrNotFound{})
		return
	}

	return
}

func (i *index2[T, TPtr]) StoreDelta(d schnittstellen.Delta[T]) (err error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	log.Log().Printf("delta: %s", d)

	if err = d.GetAdded().Each(i.storeOne); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = d.GetRemoved().Each(
		func(e T) (err error) {
			id, ok := i.Kennungen[e.String()]

			if !ok {
				err = errors.Errorf("tried to remove %s but it wasn't present", e)
				return
			}

			id.SchwanzenCount -= 1

			log.Log().Printf("new SchwanzenCount: %s -> %d", e, id.SchwanzenCount)

			i.Kennungen[e.String()] = id

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *index2[T, TPtr]) StoreMany(ks schnittstellen.SetLike[T]) (err error) {
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
	id, ok := i.Kennungen[k.String()]

	if !ok {
		id.ResetWithKennung(k)
	}

	i.hasChanges = true
	id.SchwanzenCount += 1
	id.Count += 1

	i.Kennungen[k.String()] = id

	return
}
