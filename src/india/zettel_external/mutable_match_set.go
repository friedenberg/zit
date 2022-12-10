package zettel_external

import (
	"io"

	"github.com/friedenberg/zit/src/kilo/zettel"
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

func (s MutableMatchSet) Match(z *zettel.Transacted) (err error) {
	kStored := z.Sku.Sha.String()
	kAkte := z.Objekte.Akte.String()

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
