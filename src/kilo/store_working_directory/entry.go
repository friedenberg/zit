package store_working_directory

import (
	"time"

	"github.com/friedenberg/zit/src/bravo/sha"
)

type Entry struct {
	Time time.Time
	Sha  sha.Sha
}

type EntryMap map[string]Entry

type Entries struct {
	Zettelen EntryMap
	Akten    EntryMap
}

func newEntryMap() EntryMap {
	return make(map[string]Entry)
}

func newEntries() Entries {
	return Entries{
		Zettelen: newEntryMap(),
		Akten:    newEntryMap(),
	}
}

func (em EntryMap) NormalizePath(p string) (p1 string) {
	p1 = p
	return
}

func (em EntryMap) Del(p string) {
}

func (em EntryMap) Set(p string, e Entry) {
}

func (em EntryMap) Get(p string) (e Entry, ok bool) {

	return
}
