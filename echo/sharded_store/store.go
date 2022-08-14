package sharded_store

import (
	"sync"

	"github.com/friedenberg/zit/alfa/logz"
	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/bravo/open_file_guard"
)

type Store interface {
	Storer
	Shard(id string) (s Shard, err error)
	All() (es []Entry, err error)
	Flush() error
}

type Storer interface {
	NewShard(path string, id string) (Shard, error)
}

type store struct {
	path   string
	rwLock *sync.RWMutex
	shards map[string]Shard
	Storer
}

func NewStore(path string, st Storer) (s *store, err error) {
	s = &store{
		path:   path,
		rwLock: &sync.RWMutex{},
		shards: make(map[string]Shard),
		Storer: st,
	}

	return
}

func (ss store) Shard(id string) (s Shard, err error) {
	var ok bool

	ss.rwLock.RLock()

	if s, ok = ss.shards[id]; ok {
		ss.rwLock.RUnlock()
		return
	}

	ss.rwLock.RUnlock()

	if s, err = ss.NewShard(ss.path, id); err != nil {
		err = errors.Error(err)
		return
	}

	ss.rwLock.Lock()
	ss.shards[id] = s
	ss.rwLock.Unlock()

	return
}

func (ss store) Flush() (err error) {
	for fn, s := range ss.shards {
		if err = s.Flush(); err != nil {
			err = errors.Errorf("failed to flush shard: %s: %s", fn, err)
			return
		}
	}

	return
}

func (ss store) All() (es []Entry, err error) {
	var files []string

	if files, err = open_file_guard.ReadDirNames(ss.path); err != nil {
		err = errors.Error(err)
		return
	}

	for _, fn := range files {
		var s Shard

		if s, err = ss.Shard(fn); err != nil {
			logz.Print(s)
			err = errors.Error(err)
			return
		}

		for _, e := range s.All() {
			es = append(es, e)
		}
	}

	return
}
