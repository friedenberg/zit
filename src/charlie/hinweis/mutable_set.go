package hinweis

import "github.com/friedenberg/zit/src/bravo/proto_objekte"

type MutableSet = proto_objekte.MutableSet[Hinweis, *Hinweis]

func MakeMutableSet(hs ...Hinweis) MutableSet {
	return MutableSet(proto_objekte.MakeMutableSet[Hinweis, *Hinweis](hs...))
}
