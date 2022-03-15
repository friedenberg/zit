package sharded_store

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type Sharder interface {
	LineToEntry(line string) (entry Entry, err error)
	EntryToLine(entry Entry) (line string, err error)
}

type Shard interface {
	Sharder
	Path() string
	ReadAll(*bufio.Reader) error
	Read(key string) (value string, ok bool)
	All() (a []Entry)
	Flush() error
	Set(key, value string)
	Remove(key string)
}

type shard struct {
	path       string
	age        _Age
	rwLock     *sync.RWMutex
	lines      []string
	entries    EntryMap
	hasChanges bool
	Sharder
}

func NewShard(path string, age _Age, sh Sharder) (s *shard, err error) {
	s = &shard{
		path:    path,
		age:     age,
		rwLock:  &sync.RWMutex{},
		entries: make(EntryMap),
		Sharder: sh,
	}

	var file *os.File

	if file, err = _OpenFile(s.Path(), os.O_RDONLY, 0644); err != nil {
		if os.IsNotExist(err) {
			err = nil
			return
		}

		err = _Error(err)
		return
	}

	defer _Close(file)

	var r io.ReadCloser

	if r, err = s.Reader(file); err != nil {
		err = _Error(err)
		return
	}

	defer r.Close()

	if err = s.ReadAll(bufio.NewReader(r)); err != nil {
		err = _Error(err)
		return
	}

	return
}

func (s shard) Reader(r io.ReadCloser) (ro io.ReadCloser, err error) {
	if s.age == nil {
		ro = r
		return
	}

	if ro, err = _NewObjekteReader(s.age, r); err != nil {
		err = _Error(err)
		return
	}

	return
}

func (s shard) Writer(w io.WriteCloser) (wo io.WriteCloser, err error) {
	if s.age == nil {
		wo = w
		return
	}

	if wo, err = _NewObjekteWriter(s.age, w); err != nil {
		err = _Error(err)
		return
	}

	return
}

func (s shard) Path() string {
	return s.path
}

func (s shard) ReadAll(r *bufio.Reader) (err error) {
	for {
		var l string
		l, err = r.ReadString('\n')

		if errors.Is(err, io.EOF) {
			err = nil
			break
		}

		if err != nil {
			err = _Error(err)
			return
		}

		s.lines = append(s.lines, l)
		s.processLine(l)
	}

	return
}

func (s shard) processLine(line string) (err error) {
	var entry Entry

	if entry, err = s.LineToEntry(line); err != nil {
		err = _Error(err)
		return
	}

	s.Set(entry.Key, entry.Value)

	return
}

func (s *shard) Remove(key string) {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	s.hasChanges = true
	delete(s.entries, key)
}

func (s *shard) Set(key, value string) {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	s.hasChanges = true
	s.entries[key] = value
}

func (s shard) Read(key string) (value string, ok bool) {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	if value, ok = s.entries[key]; ok {
		return
	}

	return
}

func (s shard) Flush() (err error) {
	s.rwLock.RLock()
	// maybe a mistake
	defer s.rwLock.RUnlock()

	if !s.hasChanges {
		return
	}

	log.Printf("flushing: %s", s.path)

	var file *os.File

	if file, err = _TempFile(); err != nil {
		log.Print(err)
		err = _Error(err)
		return
	}

	defer _Close(file)

	var w io.WriteCloser

	if w, err = s.Writer(file); err != nil {
		log.Print(err)
		err = _Error(err)
		return
	}

	defer w.Close()

	for k, v := range s.entries {
		var line string

		if line, err = s.EntryToLine(Entry{k, v}); err != nil {
			log.Print(err)
			err = _Error(err)
			return
		}

		if _, err = io.WriteString(w, fmt.Sprintln(line)); err != nil {
			log.Print(err)
			err = _Error(err)
			return
		}
	}

	log.Printf("renaming %s to %s", file.Name(), s.path)
	if err = os.Rename(file.Name(), s.path); err != nil {
		err = _Error(err)
		return
	}

	log.Print("done renaming")

	return
}

func (s shard) All() (a []Entry) {
	a = make([]Entry, 0, len(s.entries))

	for k, v := range s.entries {
		a = append(a, Entry{k, v})
	}

	return
}
