package kennung

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
)

func init() {
	gob.Register(&matcherEtiketten{})
}

type Matcher interface {
	ContainsMatchable(Matchable) bool
}

func MakeMatcherAlways() Matcher {
	return matcherAlways{}
}

func MakeMatcherNever() Matcher {
	return matcherNever{}
}

func MakeMatcherEtiketten(es EtikettSet) Matcher {
	return matcherEtiketten{es}
}

type matcherAlways struct{}

func (_ matcherAlways) ContainsMatchable(_ Matchable) bool {
	return true
}

type matcherNever struct{}

func (_ matcherNever) ContainsMatchable(_ Matchable) bool {
	return false
}

type matcherEtiketten struct {
	EtikettSet
}

func (f matcherEtiketten) ContainsMatchable(m Matchable) bool {
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
