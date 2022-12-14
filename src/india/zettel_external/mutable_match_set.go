package zettel_external

import (
	"io"
	"sync"

	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/kilo/zettel"
)

type MutableMatchSet struct {
	lock             *sync.RWMutex
	Original         MutableSet
	Stored           MutableSet
	Akten            MutableSet
	Matched          MutableSet
	MatchedHinweisen hinweis.MutableSet
}

func MakeMutableMatchSet(in MutableSet) (out MutableMatchSet) {
	out = MutableMatchSet{
		lock:             &sync.RWMutex{},
		Original:         in,
		Stored:           MakeMutableSetUniqueStored(),
		Akten:            MakeMutableSetUniqueAkte(),
		Matched:          MakeMutableSetUniqueFD(),
		MatchedHinweisen: hinweis.MakeMutableSet(),
	}

	in.Each(out.Stored.Add)
	in.Each(out.Akten.Add)

	return
}

func (s MutableMatchSet) Match(z *zettel.Transacted) (err error) {
	kStored := z.Sku.Sha.String()
	kAkte := z.Objekte.Akte.String()

	s.lock.RLock()
	stored, okStored := s.Stored.Get(kStored)
	akte, okAkte := s.Akten.Get(kAkte)
	okHinweis := s.MatchedHinweisen.Contains(z.Sku.Kennung)
	s.lock.RUnlock()

  //TODO-P0 figure out why matches aren't the last-most zettel
	if okStored || okAkte || okHinweis {
		s.lock.Lock()
		defer s.lock.Unlock()

		s.MatchedHinweisen.Add(z.Sku.Kennung)
		s.Stored.DelKey(kStored)
		s.Akten.DelKey(kAkte)
		s.Original.Del(stored)
		s.Original.Del(akte)

		//Only one is necessary
		s.Matched.Add(akte)

		return
	}

	err = io.EOF

	return
}
