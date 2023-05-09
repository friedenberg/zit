package kennung

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/todo"
)

func init() {
	gob.Register(&matcherEtiketten{})
}

type Matcher interface {
	ContainsMatchable(Matchable) bool
	// schnittstellen.Stringer
}

type MatcherParent interface {
	Each(schnittstellen.FuncIter[Matcher]) error
}

type MatcherParentPtr interface {
	MatcherParent
	Add(Matcher) error
}

func VisitAllMatchers(
	f schnittstellen.FuncIter[Matcher],
	matchers ...Matcher,
) (err error) {
	for _, m := range matchers {
		if err = f(m); err != nil {
			err = errors.Wrap(err)
			return
		}

		mp, ok := m.(MatcherParent)

		if !ok {
			continue
		}

		if err = mp.Each(
			func(m Matcher) (err error) {
				return VisitAllMatchers(f, m)
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
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

//      _              _
//     / \   _ __   __| |
//    / _ \ | '_ \ / _` |
//   / ___ \| | | | (_| |
//  /_/   \_\_| |_|\__,_|
//

func MakeMatcherAnd(ms ...Matcher) matcherAnd {
	return matcherAnd(ms)
}

type matcherAnd []Matcher

func (matcher *matcherAnd) Add(m Matcher) (err error) {
	*matcher = append(*matcher, m)
	return
}

func (matcher matcherAnd) ContainsMatchable(matchable Matchable) bool {
	if len(matcher) == 0 {
		return true
	}

	for _, m := range matcher {
		if !m.ContainsMatchable(matchable) {
			return false
		}
	}

	return true
}

func (matcher matcherAnd) Each(f schnittstellen.FuncIter[Matcher]) (err error) {
	for _, m := range matcher {
		if err = f(m); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

//    ___
//   / _ \ _ __
//  | | | | '__|
//  | |_| | |
//   \___/|_|
//

func MakeMatcherOr(ms ...Matcher) matcherOr {
	return matcherOr(ms)
}

type matcherOr []Matcher

func (matcher *matcherOr) Add(m Matcher) (err error) {
	*matcher = append(*matcher, m)
	return
}

func (matcher matcherOr) ContainsMatchable(matchable Matchable) bool {
	if len(matcher) == 0 {
		return true
	}

	for _, m := range matcher {
		if m.ContainsMatchable(matchable) {
			return true
		}
	}

	return false
}

func (matcher matcherOr) Each(f schnittstellen.FuncIter[Matcher]) (err error) {
	for _, m := range matcher {
		if err = f(m); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

//   _   _                  _
//  | \ | | ___  __ _  __ _| |_ ___
//  |  \| |/ _ \/ _` |/ _` | __/ _ \
//  | |\  |  __/ (_| | (_| | ||  __/
//  |_| \_|\___|\__, |\__,_|\__\___|
//              |___/

func MakeMatcherNegate(m Matcher) Matcher {
	return matcherNegate{Matcher: m}
}

type matcherNegate struct {
	Matcher
}

func (matcher matcherNegate) ContainsMatchable(matchable Matchable) bool {
	return !matcher.Matcher.ContainsMatchable(matchable)
}

func (matcher matcherNegate) Each(f schnittstellen.FuncIter[Matcher]) error {
	return f(matcher.Matcher)
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
