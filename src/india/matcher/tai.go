package matcher

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/kennung"
)

type Tai kennung.Tai

func (t Tai) ContainsMatchableExactly(m Matchable) bool {
	return false
}

func (t Tai) ContainsMatchable(m Matchable) bool {
	errors.TodoP1("add GetTai to matchable")
	return false
}
