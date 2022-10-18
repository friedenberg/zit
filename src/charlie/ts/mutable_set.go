package ts

import "github.com/friedenberg/zit/src/bravo/proto_objekte"

type MutableSet = proto_objekte.MutableSet[Time, *Time]

func MakeMutableSet(hs ...Time) MutableSet {
	return MutableSet(proto_objekte.MakeMutableSet[Time, *Time](hs...))
}
