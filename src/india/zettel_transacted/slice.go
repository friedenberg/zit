package zettel_transacted

import (
	"encoding/json"
	"sort"
)

type Slice struct {
	innerSlice []Zettel
}

func MakeSliceTransacted() Slice {
	return Slice{
		innerSlice: make([]Zettel, 0),
	}
}

func (s Slice) Len() int {
	return len(s.innerSlice)
}

func (s *Slice) Append(tz Zettel) {
	s.innerSlice = append(s.innerSlice, tz)
}

func (s Slice) Get(i int) Zettel {
	return s.innerSlice[i]
}

func (s *Slice) Sort(f func(int, int) bool) {
	sort.Slice(s.innerSlice, f)
}

func (s Slice) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.innerSlice)
}

func (s *Slice) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &s.innerSlice)
}
