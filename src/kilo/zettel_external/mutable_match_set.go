package zettel_external

import (
	"sync"

	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/echo/hinweis"
	"github.com/friedenberg/zit/src/echo/ts"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

type MutableMatchSet struct {
	lock                      *sync.RWMutex
	Original                  MutableSet
	Stored                    MutableSet
	Akten                     MutableSet
	Matched                   MutableSet
	MatchedHinweisen          hinweis.MutableSet
	MatchedHinweisenSchwanzen map[hinweis.Hinweis]ts.Time
}

func MakeMutableMatchSet(in MutableSet) (out MutableMatchSet) {
	out = MutableMatchSet{
		lock:                      &sync.RWMutex{},
		Original:                  in,
		Stored:                    MakeMutableSetUniqueStored(),
		Akten:                     MakeMutableSetUniqueAkte(),
		Matched:                   MakeMutableSetUniqueFD(),
		MatchedHinweisen:          hinweis.MakeMutableSet(),
		MatchedHinweisenSchwanzen: make(map[hinweis.Hinweis]ts.Time),
	}

	in.Each(out.Stored.Add)
	in.Each(out.Akten.Add)

	return
}

func (s MutableMatchSet) Match(z *zettel.Transacted) (err error) {
	kStored := z.Sku.ObjekteSha.String()
	kAkte := z.Objekte.Akte.String()

	s.lock.RLock()
	stored, okStored := s.Stored.Get(kStored)
	akte, okAkte := s.Akten.Get(kAkte)
	okHinweis := s.MatchedHinweisen.Contains(z.Sku.Kennung)

	okSchwanz := false
	schwanz, _ := s.MatchedHinweisenSchwanzen[z.Sku.Kennung]

	if schwanz.Less(z.Sku.Schwanz) {
		okSchwanz = true
	}

	s.lock.RUnlock()

	// This function gets called out of order for all zettels because it is
	// parallelized. The only case this does not correctly handle is if the akte
	// is mutated or removed at some point in a zettel's history. Then, when
	// reading verzeichnisse, the _latest_ (highest Schwanz) zettel may pass
	// through this function _before_ the function has matched on a historical
	// akte or stored sha. In that case, the zettel would accidentally be
	// reverted.
	// TODO-P2 solve for the above
	if okStored || okAkte || (okHinweis && okSchwanz) {
		s.lock.Lock()
		defer s.lock.Unlock()

		s.MatchedHinweisen.Add(z.Sku.Kennung)
		s.MatchedHinweisenSchwanzen[z.Sku.Kennung] = z.Sku.Schwanz
		s.Stored.DelKey(kStored)
		s.Akten.DelKey(kAkte)
		s.Original.Del(stored)
		s.Original.Del(akte)

		//Only one is necessary
		s.Matched.Add(akte)

		return
	}

	err = collections.ErrStopIteration

	return
}