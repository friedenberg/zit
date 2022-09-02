package collections

import (
	"encoding/json"
	"sort"

	"github.com/friedenberg/zit/zettel_transacted"
)

type SliceTransacted struct {
	innerSlice []zettel_transacted.Transacted
}

func MakeSliceTransacted() SliceTransacted {
	return SliceTransacted{
		innerSlice: make([]zettel_transacted.Transacted, 0),
	}
}

func (s SliceTransacted) Len() int {
	return len(s.innerSlice)
}

func (s *SliceTransacted) Append(tz zettel_transacted.Transacted) {
	s.innerSlice = append(s.innerSlice, tz)
}

func (s SliceTransacted) Get(i int) zettel_transacted.Transacted {
	return s.innerSlice[i]
}

func (s *SliceTransacted) Sort(f func(int, int) bool) {
	sort.Slice(s.innerSlice, f)
}

func (s SliceTransacted) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.innerSlice)
}

func (s *SliceTransacted) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &s.innerSlice)
}
