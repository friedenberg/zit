package collections

import (
	"encoding/json"
	"sort"

	zettel_stored "github.com/friedenberg/zit/src/golf/zettel_stored"
)

type SliceTransacted struct {
	innerSlice []zettel_stored.Transacted
}

func MakeSliceTransacted() SliceTransacted {
	return SliceTransacted{
		innerSlice: make([]zettel_stored.Transacted, 0),
	}
}

func (s SliceTransacted) Len() int {
	return len(s.innerSlice)
}

func (s *SliceTransacted) Append(tz zettel_stored.Transacted) {
	s.innerSlice = append(s.innerSlice, tz)
}

func (s SliceTransacted) Get(i int) zettel_stored.Transacted {
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
