package zettel_external

import (
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type MutableMatchSet struct {
	lock                      *sync.RWMutex
	Original                  MutableSet
	Stored                    MutableSet
	Akten                     MutableSet
	Matched                   MutableSet
	MatchedHinweisen          schnittstellen.MutableSetLike[kennung.Kennung]
	MatchedHinweisenSchwanzen map[string]kennung.Tai
}

func MakeMutableMatchSet(in MutableSet) (out MutableMatchSet) {
	out = MutableMatchSet{
		lock:     &sync.RWMutex{},
		Original: in,
		Stored:   MakeMutableSetUniqueStored(),
		Akten:    MakeMutableSetUniqueAkte(),
		Matched:  MakeMutableSetUniqueFD(),
		MatchedHinweisen: collections_value.MakeMutableValueSet[kennung.Kennung](
			nil,
		),
		MatchedHinweisenSchwanzen: make(map[string]kennung.Tai),
	}

	in.Each(out.Stored.Add)
	in.Each(out.Akten.Add)

	return
}

func (s MutableMatchSet) Match(z sku.SkuLikePtr) (err error) {
	kStored := z.GetObjekteSha().String()
	kAkte := z.GetAkteSha().String()

	s.lock.RLock()
	stored, okStored := s.Stored.Get(kStored)
	akte, okAkte := s.Akten.Get(kAkte)
	k := z.GetKennungLike()
	okHinweis := s.MatchedHinweisen.Contains(z.GetKennungLike())

	okSchwanz := false
	schwanz, _ := s.MatchedHinweisenSchwanzen[k.String()]

	if schwanz.Less(z.GetTai()) {
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
	errors.TodoP2("fix history erasure on zettel match")
	if okStored || okAkte || (okHinweis && okSchwanz) {
		s.lock.Lock()
		defer s.lock.Unlock()

		s.MatchedHinweisen.Add(k)
		s.MatchedHinweisenSchwanzen[k.String()] = z.GetTai()
		s.Stored.DelKey(kStored)
		s.Akten.DelKey(kAkte)
		s.Original.Del(stored)
		s.Original.Del(akte)

		// Only one is necessary
		s.Matched.Add(akte)

		return
	}

	err = collections.MakeErrStopIteration()

	return
}
