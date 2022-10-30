package zettel_external

import (
	"io"

	collections "github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
)

type MutableMatchSet struct {
	Original MutableSet
	Stored   MutableSet
	Akten    MutableSet
}

func MakeMutableMatchSet(in MutableSet) (out MutableMatchSet) {
	out = MutableMatchSet{
		Original: in,
		Stored:   MakeMutableSetUniqueStored(),
		Akten:    MakeMutableSetUniqueAkte(),
	}

	in.Each(out.Stored.Add)
	in.Each(out.Akten.Add)

	return
}

func (s MutableMatchSet) WriterZettelNamed() collections.WriterFunc[*zettel_named.Zettel] {
	return func(z *zettel_named.Zettel) (err error) {
		kStored := z.Stored.Sha.String()
		kAkte := z.Stored.Zettel.Akte.String()

		stored, okStored := s.Stored.Get(kStored)
		akte, okAkte := s.Akten.Get(kAkte)

		if okStored || okAkte {
			s.Stored.DelKey(kStored)
			s.Akten.DelKey(kAkte)
			s.Original.Del(stored)
			s.Original.Del(akte)
			return
		}

		err = io.EOF

		return
	}
}
