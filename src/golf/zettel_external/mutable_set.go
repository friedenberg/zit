package zettel_external

import "github.com/friedenberg/zit/src/bravo/proto_objekte"

type MutableSet struct {
	proto_objekte.MutableSetGeneric[*Zettel]
}

func MakeMutableSet(zs ...*Zettel) MutableSet {
	kf := func(z *Zettel) string {
		return z.String()
	}

	return MutableSet{
		MutableSetGeneric: proto_objekte.MakeMutableSetGeneric[*Zettel](kf, zs...),
	}
}
