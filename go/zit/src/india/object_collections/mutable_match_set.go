package object_collections

import (
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type MutableMatchSet struct {
	lock                      *sync.RWMutex
	Original                  MutableSet
	Stored                    MutableSet
	Akten                     MutableSet
	Matched                   MutableSet
	MatchedHinweisen          interfaces.MutableSetLike[ids.IdLike]
	MatchedHinweisenSchwanzen map[string]ids.Tai
}

func MakeMutableMatchSet(in MutableSet) (out MutableMatchSet) {
	out = MutableMatchSet{
		lock:     &sync.RWMutex{},
		Original: in,
		Stored:   MakeMutableSetUniqueStored(),
		Akten:    MakeMutableSetUniqueAkte(),
		Matched:  MakeMutableSetUniqueFD(),
		MatchedHinweisen: collections_value.MakeMutableValueSet[ids.IdLike](
			nil,
		),
		MatchedHinweisenSchwanzen: make(map[string]ids.Tai),
	}

	in.Each(out.Stored.Add)
	in.Each(out.Akten.Add)

	return
}

func (s MutableMatchSet) Match(z *sku.Transacted) (err error) {
	kStored := z.GetObjectSha().String()
	kAkte := z.GetBlobSha().String()

	s.lock.RLock()
	stored, okStored := s.Stored.Get(kStored)
	akte, okAkte := s.Akten.Get(kAkte)
	k := z.GetObjectId()
	okHinweis := s.MatchedHinweisen.Contains(z.GetObjectId())

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