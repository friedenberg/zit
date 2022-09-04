package zettel_transacted

import (
	"sort"
)

type Slice []Zettel

func MakeSlice(c int) Slice {
	return make([]Zettel, 0, c)
}

func (s Slice) Len() int {
	return len(s)
}

func (s *Slice) Append(tz Zettel) {
	*s = append(*s, tz)
}

func (s Slice) Get(i int) Zettel {
	return s[i]
}

func (s *Slice) Sort(f func(int, int) bool) {
	sort.Slice(*s, f)
}
