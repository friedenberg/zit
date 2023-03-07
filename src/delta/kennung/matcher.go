package kennung

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
)

type Matcher interface {
	ContainsMatchable(Matchable) bool
}

func MakeMatcherAlways() Matcher {
	return matcherAlways{}
}

func MakeMatcherNever() Matcher {
	return matcherNever{}
}

type matcherAlways struct{}

func (_ matcherAlways) ContainsMatchable(_ Matchable) bool {
	return true
}

type matcherNever struct{}

func (_ matcherNever) ContainsMatchable(_ Matchable) bool {
	return false
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
