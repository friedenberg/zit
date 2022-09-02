package zettel_transacted

import (
	"encoding/json"
	"sort"
)

type SliceTransacted struct {
	innerSlice []Transacted
}

func MakeSliceTransacted() SliceTransacted {
	return SliceTransacted{
		innerSlice: make([]Transacted, 0),
	}
}

func (s SliceTransacted) Len() int {
	return len(s.innerSlice)
}

func (s *SliceTransacted) Append(tz Transacted) {
	s.innerSlice = append(s.innerSlice, tz)
}

func (s SliceTransacted) Get(i int) Transacted {
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
