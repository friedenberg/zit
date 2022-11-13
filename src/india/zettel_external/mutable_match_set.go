package zettel_external

import (
	"io"

	"github.com/friedenberg/zit/src/hotel/zettel_named"
)

type MutableMatchSet struct {
	Original MutableSet
	Stored   MutableSet
	Akten    MutableSet
	Matched  MutableSet
}

func MakeMutableMatchSet(in MutableSet) (out MutableMatchSet) {
	out = MutableMatchSet{
		Original: in,
		Stored:   MakeMutableSetUniqueStored(),
		Akten:    MakeMutableSetUniqueAkte(),
		Matched:  MakeMutableSetUniqueStored(),
	}

	in.Each(out.Stored.Add)
	in.Each(out.Akten.Add)

	return
}

func (s MutableMatchSet) Match(z *zettel_named.Zettel) (err error) {
	kStored := z.Stored.Sha.String()
	kAkte := z.Stored.Zettel.Akte.String()

	stored, okStored := s.Stored.Get(kStored)
	akte, okAkte := s.Akten.Get(kAkte)

	if okStored || okAkte {
		s.Stored.DelKey(kStored)
		s.Akten.DelKey(kAkte)
		s.Original.Del(stored)
		s.Original.Del(akte)

		//These two should be redundant
		s.Matched.Add(akte)
		s.Matched.Add(stored)
		return
	}

	err = io.EOF

	return
}
