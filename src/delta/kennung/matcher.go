package kennung

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/todo"
)

func init() {
	gob.Register(&matcherEtiketten{})
}

type Matcher interface {
	ContainsMatchable(Matchable) bool
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

//   _   _                  _
//  | \ | | ___  __ _  __ _| |_ ___
//  |  \| |/ _ \/ _` |/ _` | __/ _ \
//  | |\  |  __/ (_| | (_| | ||  __/
//  |_| \_|\___|\__, |\__,_|\__\___|
//              |___/

func MakeMatcherNegate() Matcher {
	return matcherNegate{}
}

type matcherNegate struct {
	Matcher
}

func (matcher matcherNegate) ContainsMatchable(matchable Matchable) bool {
	return !matcher.Matcher.ContainsMatchable(matchable)
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

func (f matcherEtiketten) ContainsMatchable(m Matchable) (ok bool) {
	todo.Optimize()
	ok = iter.AnyOrFalseEmpty[Etikett](f, m.GetEtikettenExpanded().Contains)
	return
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

//   _____ _   _ _        _   _
//  | ____| |_(_) | _____| |_| |_
//  |  _| | __| | |/ / _ \ __| __|
//  | |___| |_| |   <  __/ |_| |_
//  |_____|\__|_|_|\_\___|\__|\__|
//

func MakeMatcherEtikett(e Etikett) Matcher {
	return matcherEtikett{Etikett: e}
}

type matcherEtikett struct {
	Etikett
}

func (e matcherEtikett) ContainsMatchable(m Matchable) bool {
	es := m.GetEtiketten()

	if es.Contains(e.Etikett) {
		return true
	}

	e1, ok := m.GetIdLike().(Etikett)

	if ok && Contains(e1, e.Etikett) {
		return true
	}

	return false
}
