package sharded_store

import (
	"log"
	"sync"
)

type StoreImmutable interface {
	StorerImmutable
	ShardImmutable(id string) (s ShardImmutable, err error)
	All() (es []Entry, err error)
}

type StorerImmutable interface {
	NewShardImmutable(path string, id string) (ShardImmutable, error)
}

type storeImmutable struct {
	path   string
	rwLock *sync.RWMutex
	shards map[string]ShardImmutable
	StorerImmutable
}

func NewStoreImmutable(path string, st StorerImmutable) (s *storeImmutable, err error) {
	s = &storeImmutable{
		path:            path,
		rwLock:          &sync.RWMutex{},
		shards:          make(map[string]ShardImmutable),
		StorerImmutable: st,
	}

	return
}

func (ss storeImmutable) ShardImmutable(id string) (s ShardImmutable, err error) {
	var ok bool

	ss.rwLock.RLock()

	if s, ok = ss.shards[id]; ok {
		ss.rwLock.RUnlock()
		return
	}

	ss.rwLock.RUnlock()

	if s, err = ss.NewShardImmutable(ss.path, id); err != nil {
		err = _Error(err)
		return
	}

	ss.rwLock.Lock()
	ss.shards[id] = s
	ss.rwLock.Unlock()

	return
}

func (ss storeImmutable) All() (es []Entry, err error) {
	var files []string

	if files, err = _ReadDirNames(ss.path); err != nil {
		err = _Error(err)
		return
	}

	for _, fn := range files {
		var s ShardImmutable

		if s, err = ss.ShardImmutable(fn); err != nil {
			log.Print(s)
			err = _Error(err)
			return
		}

		for _, e := range s.All() {
			es = append(es, e)
		}
	}

	return
}
