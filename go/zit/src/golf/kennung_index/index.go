package kennung_index

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"io"
	"strconv"
	"sync"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/src/delta/sha"
	"code.linenisgreat.com/zit/src/echo/kennung"
)

type index2[
	T kennung.KennungLike[T],
	TPtr kennung.KennungLikePtr[T],
] struct {
	path            string
	vf              schnittstellen.VerzeichnisseFactory
	readOnce        *sync.Once
	hasChanges      bool
	lock            *sync.RWMutex
	IntsToKennungen map[int]TPtr
	Kennungen       map[string]*kennung.IndexedLike
}

func MakeIndex2[
	T kennung.KennungLike[T],
	TPtr kennung.KennungLikePtr[T],
](
	vf schnittstellen.VerzeichnisseFactory,
	path string,
) (i *index2[T, TPtr]) {
	i = &index2[T, TPtr]{
		path:            path,
		vf:              vf,
		readOnce:        &sync.Once{},
		lock:            &sync.RWMutex{},
		IntsToKennungen: make(map[int]TPtr),
		Kennungen:       make(map[string]*kennung.IndexedLike),
	}

	return
}

func (i *index2[T, TPtr]) HasChanges() bool {
	i.lock.RLock()
	defer i.lock.RUnlock()

	return i.hasChanges
}

func (i *index2[T, TPtr]) Reset() error {
	i.lock.Lock()
	defer i.lock.Unlock()

	i.Kennungen = make(map[string]*kennung.IndexedLike)
	i.IntsToKennungen = make(map[int]TPtr)
	i.readOnce = &sync.Once{}
	i.hasChanges = false

	return nil
}

func (ei *index2[T, TPtr]) Flush() (err error) {
	return ei.WriteIfNecessary()
}

func (ei *index2[T, TPtr]) WriteIfNecessary() (err error) {
	if !ei.HasChanges() {
		ui.Log().Printf("%s does not have changes", ei.path)
		return
	}

	ui.Log().Printf("%s has changes", ei.path)

	var wc schnittstellen.ShaWriteCloser

	if wc, err = ei.vf.WriteCloserVerzeichnisse(ei.path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, wc)

	if _, err = ei.WriteTo(wc); err != nil {
		err = errors.Wrap(err)
		return
	}

	ui.Log().Printf("%s done writing changes", ei.path)

	return
}

func (i *index2[T, TPtr]) WriteTo(w1 io.Writer) (n int64, err error) {
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

func (ei *index2[T, TPtr]) ReadIfNecessary() (err error) {
	ei.readOnce.Do(func() { err = ei.read() })
	return
}

func (ei *index2[T, TPtr]) read() (err error) {
	var rc io.ReadCloser

	if rc, err = ei.vf.ReadCloserVerzeichnisse(ei.path); err != nil {
		if errors.IsNotExist(err) {
			err = nil
			rc = sha.MakeReadCloser(bytes.NewBuffer(nil))
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	defer errors.DeferredCloser(&err, rc)

	if _, err = ei.ReadFrom(rc); err != nil {
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

	return
}

func (i *index2[T, TPtr]) Each(
	f schnittstellen.FuncIter[kennung.IndexedLike],
) (err error) {
	if err = i.ReadIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, id := range i.Kennungen {
		if err = f(*id); err != nil {
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

func (i *index2[T, TPtr]) EachSchwanzen(
	f schnittstellen.FuncIter[*kennung.IndexedLike],
) (err error) {
	if err = i.ReadIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

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

func (i *index2[T, TPtr]) GetAll() (out []kennung.Kennung, err error) {
	if err = i.ReadIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	out = make([]kennung.Kennung, 0, len(i.Kennungen))

	for _, ki := range i.Kennungen {
		out = append(out, ki.GetKennung())
	}

	return
}

func (i *index2[T, TPtr]) GetInt(in int) (id T, err error) {
	if err = i.ReadIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	i.lock.RLock()
	defer i.lock.RUnlock()

	ok := false
	var id1 TPtr
	id1, ok = i.IntsToKennungen[in]

	if !ok {
		err = collections.MakeErrNotFoundString(strconv.Itoa(in))
		return
	}

	id = *id1

	return
}

func (i *index2[T, TPtr]) Get(
	k TPtr,
) (id *kennung.IndexedLike, err error) {
	if err = i.ReadIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	i.lock.RLock()
	defer i.lock.RUnlock()

	ok := false
	id, ok = i.Kennungen[k.String()]

	if !ok {
		err = collections.MakeErrNotFound(k)
		return
	}

	return
}

func (i *index2[T, TPtr]) StoreDelta(d schnittstellen.Delta[T]) (err error) {
	if err = i.ReadIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	i.lock.Lock()
	defer i.lock.Unlock()

	ui.Log().Printf("delta: %s", d)

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

			ui.Log().Printf("new SchwanzenCount: %s -> %d", e, id.SchwanzenCount)

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
	if err = i.ReadIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	i.lock.Lock()
	defer i.lock.Unlock()

	return ks.Each(i.storeOne)
}

func (i *index2[T, TPtr]) StoreOne(k T) (err error) {
	if err = i.ReadIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

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
		id = &kennung.IndexedLike{}
		id.ResetWithKennung(k)
	}

	i.hasChanges = true
	id.SchwanzenCount += 1
	id.Count += 1

	ui.Log().Printf("new SchwanzenCount: %s -> %d", k, id.SchwanzenCount)

	i.Kennungen[k.String()] = id

	return
}
