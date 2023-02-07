package kennung

import "github.com/friedenberg/zit/src/delta/gattungen"

type MetaSet struct {
	Gattung gattungen.MutableSet
	IdSet   Set
}

func MakeMetaSet() *MetaSet {
	return &MetaSet{
		Gattung: gattungen.MakeMutableSet(),
		IdSet:   MakeSet(),
	}
}

func (ms *MetaSet) Set(v string) (err error) {
	return
}
