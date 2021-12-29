package sharded_store

import (
	"bufio"
	"errors"
	"io"
	"os"
)

type ShardImmutable interface {
	Sharder
	Path() string
	ReadAll(*bufio.Reader) error
	Read(key string) (value string, ok bool)
	All() (a []Entry)
}

type Sharder interface {
	LineToEntry(line string) (entry Entry, err error)
	EntryToLine(entry Entry) (line string, err error)
}

type shardImmutable struct {
	path       string
	age        _Age
	lines      []string
	entries    EntryMap
	hasChanges bool
	Sharder
}

func NewShardImmutable(path string, age _Age, sh Sharder) (s *shardImmutable, err error) {
	s = &shardImmutable{
		path:    path,
		age:     age,
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

func (s shardImmutable) Reader(r io.ReadCloser) (ro io.ReadCloser, err error) {
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

func (s shardImmutable) Writer(w io.WriteCloser) (wo io.WriteCloser, err error) {
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

func (s shardImmutable) Path() string {
	return s.path
}

func (s shardImmutable) ReadAll(r *bufio.Reader) (err error) {
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

func (s shardImmutable) processLine(line string) (err error) {
	var entry Entry

	if entry, err = s.LineToEntry(line); err != nil {
		err = _Error(err)
		return
	}

	s.entries[entry.Key] = entry.Value

	return
}

func (s shardImmutable) Read(key string) (value string, ok bool) {
	if value, ok = s.entries[key]; ok {
		return
	}

	return
}

func (s shardImmutable) All() (a []Entry) {
	a = make([]Entry, 0, len(s.entries))

	for k, v := range s.entries {
		a = append(a, Entry{k, v})
	}

	return
}
