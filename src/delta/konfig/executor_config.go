package konfig

import (
	"github.com/friedenberg/zit/src/bravo/collections"
)

type Executor struct {
	collections.StringValue
}

func (s *Executor) Merge(s2 *Executor) {
	if !s2.IsEmpty() {
		s.Set(s2.String())
	}
}
