package typ

import "github.com/friedenberg/zit/src/bravo/proto_objekte"

type MutableSet = proto_objekte.MutableSet[Typ, *Typ]

func MakeMutableSet(hs ...Typ) MutableSet {
	return MutableSet(proto_objekte.MakeMutableSet[Typ, *Typ](hs...))
}
