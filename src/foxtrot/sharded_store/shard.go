package sharded_store

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/charlie/open_file_guard"
	"github.com/friedenberg/zit/src/delta/age"
	age_io "github.com/friedenberg/zit/src/echo/age_io"
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
	age        age.Age
	rwLock     *sync.RWMutex
	lines      []string
	entries    EntryMap
	hasChanges bool
	Sharder
}

func NewShard(path string, age age.Age, sh Sharder) (s *shard, err error) {
	s = &shard{
		path:    path,
		age:     age,
		rwLock:  &sync.RWMutex{},
		entries: make(EntryMap),
		Sharder: sh,
	}

	var file *os.File

	if file, err = open_file_guard.OpenFile(s.Path(), os.O_RDONLY, 0644); err != nil {
		if os.IsNotExist(err) {
			err = nil
			return
		}

		err = errors.Error(err)
		return
	}

	defer open_file_guard.Close(file)

	var r io.ReadCloser

	if r, err = s.Reader(file); err != nil {
		err = errors.Error(err)
		return
	}

	defer r.Close()

	if err = s.ReadAll(bufio.NewReader(r)); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (s shard) Reader(r io.ReadCloser) (ro io.ReadCloser, err error) {
	if s.age == nil {
		ro = r
		return
	}

	o := age_io.ReadOptions{
		Age:    s.age,
		Reader: r,
	}

	if ro, err = age_io.NewReader(o); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (s shard) Writer(w io.WriteCloser) (wo io.WriteCloser, err error) {
	if s.age == nil {
		wo = w
		return
	}

	o := age_io.WriteOptions{
		Age:    s.age,
		Writer: w,
	}

	if wo, err = age_io.NewWriter(o); err != nil {
		err = errors.Error(err)
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
			err = errors.Error(err)
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
		err = errors.Error(err)
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

	logz.Printf("flushing: %s", s.path)

	var file *os.File

	if file, err = open_file_guard.TempFile(); err != nil {
		logz.Print(err)
		err = errors.Error(err)
		return
	}

	defer open_file_guard.Close(file)

	var w io.WriteCloser

	if w, err = s.Writer(file); err != nil {
		logz.Print(err)
		err = errors.Error(err)
		return
	}

	defer w.Close()

	for k, v := range s.entries {
		var line string

		if line, err = s.EntryToLine(Entry{k, v}); err != nil {
			logz.Print(err)
			err = errors.Error(err)
			return
		}

		if _, err = io.WriteString(w, fmt.Sprintln(line)); err != nil {
			logz.Print(err)
			err = errors.Error(err)
			return
		}
	}

	//TODO-research should the file be closed before being renamed???
	logz.Printf("renaming %s to %s", file.Name(), s.path)
	if err = os.Rename(file.Name(), s.path); err != nil {
		err = errors.Error(err)
		return
	}

	logz.Print("done renaming")

	return
}

func (s shard) All() (a []Entry) {
	a = make([]Entry, 0, len(s.entries))

	for k, v := range s.entries {
		a = append(a, Entry{k, v})
	}

	return
}
