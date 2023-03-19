package kennung

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/delta/gattungen"
)

func init() {
	gob.Register(&matcherEtiketten{})
}

type Matcher interface {
	ContainsMatchable(Matchable) bool
}

type MatcherMutable interface {
	Matcher
	Add(schnittstellen.ValueLike, Sigil) error
}

type GattungMatcherMap map[gattung.Gattung]Matcher

func MakeGattungMatcherMap(gs gattungen.Set, matcher Matcher) GattungMatcherMap {
	m := make(GattungMatcherMap, gs.Len())

	gs.Each(
		func(g gattung.Gattung) (err error) {
			m[g] = matcher
			return
		},
	)

	return m
}

//      _    _
//     / \  | |_      ____ _ _   _ ___
//    / _ \ | \ \ /\ / / _` | | | / __|
//   / ___ \| |\ V  V / (_| | |_| \__ \
//  /_/   \_\_| \_/\_/ \__,_|\__, |___/
//                           |___/

func MakeMatcherAlways() Matcher {
	return matcherAlways{}
}

type matcherAlways struct{}

func (_ matcherAlways) ContainsMatchable(_ Matchable) bool {
	return true
}

//   _   _
//  | \ | | _____   _____ _ __
//  |  \| |/ _ \ \ / / _ \ '__|
//  | |\  |  __/\ V /  __/ |
//  |_| \_|\___| \_/ \___|_|
//

func MakeMatcherNever() Matcher {
	return matcherNever{}
}

type matcherNever struct{}

func (_ matcherNever) ContainsMatchable(_ Matchable) bool {
	return false
}

//   _____ _   _ _        _   _
//  | ____| |_(_) | _____| |_| |_ ___ _ __
//  |  _| | __| | |/ / _ \ __| __/ _ \ '_ \
//  | |___| |_| |   <  __/ |_| ||  __/ | | |
//  |_____|\__|_|_|\_\___|\__|\__\___|_| |_|
//

func MakeMatcherEtiketten(es EtikettSet) Matcher {
	return matcherEtiketten{es.MutableClone()}
}

type matcherEtiketten struct {
	schnittstellen.MutableSet[Etikett]
}

func (f matcherEtiketten) ContainsMatchable(m Matchable) bool {
	todo.Optimize()
	return iter.AnyOrFalseEmpty[Etikett](f, m.GetEtiketten().Contains)
}

func MakeMatcherFuncIter[T Matchable](m Matcher) schnittstellen.FuncIter[T] {
	return func(e T) (err error) {
		if !m.ContainsMatchable(e) {
			err = iter.MakeErrStopIteration()
			return
		}

		return
	}
}
