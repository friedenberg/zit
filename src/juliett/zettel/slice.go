package zettel

import (
	"sort"
)

type Slice []Transacted

func MakeSlice(c int) Slice {
	return make([]Transacted, 0, c)
}

func (s Slice) Len() int {
	return len(s)
}

func (s *Slice) Append(tz Transacted) {
	*s = append(*s, tz)
}

func (s Slice) Get(i int) Transacted {
	return s[i]
}

func (s *Slice) Sort(f func(int, int) bool) {
	sort.Slice(*s, f)
}
