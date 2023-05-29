package etiketten_index

import (
	"bufio"
	"encoding/gob"
	"io"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type index struct {
	lock      *sync.Mutex
	Etiketten map[kennung.Etikett]indexed
}

func makeIndex() (i *index) {
	i = &index{
		lock:      &sync.Mutex{},
		Etiketten: make(map[kennung.Etikett]indexed),
	}

	return
}

func (i index) WriteTo(w1 io.Writer) (n int64, err error) {
	w := bufio.NewWriter(w1)
	defer errors.DeferredFlusher(&err, w)

	i.lock.Lock()
	defer i.lock.Unlock()

	enc := gob.NewEncoder(w)

	if err = enc.Encode(i.Etiketten); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *index) ReadFrom(r1 io.Reader) (n int64, err error) {
	r := bufio.NewReader(r1)

	dec := gob.NewDecoder(r)

	i.lock.Lock()
	defer i.lock.Unlock()

	if err = dec.Decode(&i.Etiketten); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i index) ExpandEtikett(k kennung.Etikett) (id Indexed, ok bool) {
	i.lock.Lock()
	defer i.lock.Unlock()

	id, ok = i.Etiketten[k]
	return
}

func (i *index) StoreEtiketten(es kennung.EtikettSet) (err error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	if err = es.Each(i.storeEtikett); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *index) StoreEtikett(k kennung.Etikett) (err error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	return i.storeEtikett(k)
}

func (i *index) storeEtikett(k kennung.Etikett) (err error) {
	if _, ok := i.Etiketten[k]; ok {
		return
	}

	id := indexed{}
	id.ResetWithEtikett(k)
	i.Etiketten[k] = id
	return
}
