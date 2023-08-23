package matcher

import "github.com/friedenberg/zit/src/echo/kennung"

type Time kennung.Time

func (t Time) ContainsMatchable(m Matchable) bool {
	return false
}
