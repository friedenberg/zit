package typen_index

import (
	"bufio"
	"encoding/gob"
	"io"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type index struct {
	didRead    bool
	hasChanges bool
	lock       *sync.Mutex
	Typen      map[kennung.Typ]indexed
}

func makeIndex() (i *index) {
	i = &index{
		lock:  &sync.Mutex{},
		Typen: make(map[kennung.Typ]indexed),
	}

	return
}

func (i *index) DidRead() bool {
	i.lock.Lock()
	defer i.lock.Unlock()

	return i.didRead
}

func (i *index) HasChanges() bool {
	i.lock.Lock()
	defer i.lock.Unlock()

	return i.hasChanges
}

func (i *index) Reset() error {
	i.lock.Lock()
	defer i.lock.Unlock()

	i.Typen = make(map[kennung.Typ]indexed)

	return nil
}

func (i index) WriteTo(w1 io.Writer) (n int64, err error) {
	w := bufio.NewWriter(w1)
	defer errors.DeferredFlusher(&err, w)

	i.lock.Lock()
	defer i.lock.Unlock()

	enc := gob.NewEncoder(w)

	if err = enc.Encode(i.Typen); err != nil {
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

	if err = dec.Decode(&i.Typen); err != nil {
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

func (i index) ExpandTyp(k kennung.Typ) (id Indexed, ok bool) {
	i.lock.Lock()
	defer i.lock.Unlock()

	id, ok = i.Typen[k]
	return
}

func (i *index) StoreTyp(k kennung.Typ) (err error) {
	if k.IsEmpty() {
		return
	}

	i.lock.Lock()
	defer i.lock.Unlock()

	return i.storeTyp(k)
}

func (i *index) storeTyp(k kennung.Typ) (err error) {
	if _, ok := i.Typen[k]; ok {
		return
	}

	i.hasChanges = true

	id := indexed{}
	id.ResetWithTyp(k)
	i.Typen[k] = id
	return
}
