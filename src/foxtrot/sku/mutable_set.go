package sku

import (
	"bytes"
	"encoding/gob"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/collections"
)

type codable struct {
	Objekten map[string][]SkuLike
	Count    int
}

type MutableSet struct {
	lock    *sync.RWMutex
	codable codable
}

func MakeMutableSet() MutableSet {
	return MutableSet{
		lock: &sync.RWMutex{},
		codable: codable{
			Objekten: make(map[string][]SkuLike),
		},
	}
}

func (os *MutableSet) Len() int {
	return os.codable.Count
}

func (os *MutableSet) Add(o SkuLike) (i int) {
	os.codable.Count++
	k := o.GetKey()

	os.lock.RLock()
	s, _ := os.codable.Objekten[k]
	os.lock.RUnlock()

	i = len(s)
	s = append(s, o)

	os.lock.Lock()
	os.codable.Objekten[k] = s
	os.lock.Unlock()

	return
}

func (os MutableSet) Get(k string) []SkuLike {
	os.lock.RLock()
	defer os.lock.RUnlock()

	return os.codable.Objekten[k]
}

func (os MutableSet) Each(
	w collections.WriterFunc[SkuLike],
) (err error) {
	os.lock.RLock()
	defer os.lock.RUnlock()

	for _, oss := range os.codable.Objekten {
		for _, o := range oss {
			if err = w(o); err != nil {
				switch {
				case collections.IsStopIteration(err):
					err = nil
					return

				default:
					err = errors.Wrap(err)
					return
				}
			}
		}
	}

	return
}

func (os MutableSet) GobEncode() (bs []byte, err error) {
	b := bytes.NewBuffer(bs)
	enc := gob.NewEncoder(b)

	if err = enc.Encode(os.codable); err != nil {
		err = errors.Wrap(err)
		return
	}

	bs = b.Bytes()

	return
}

func (os *MutableSet) GobDecode(bs []byte) (err error) {
	b := bytes.NewBuffer(bs)
	dec := gob.NewDecoder(b)

	if err = dec.Decode(&os.codable); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}