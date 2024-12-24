package sha_probe_index

import (
	"bytes"
	"slices"
	"sort"
)

type addedSlice []*row

func (a *addedSlice) GetSortable() sort.Interface {
	return a
}

func (a *addedSlice) Len() int {
	return len(*a)
}

func (s *addedSlice) Less(i, j int) bool {
	a, b := (*s)[i], (*s)[j]

	cmp := bytes.Compare(a.left.GetShaBytes(), b.left.GetShaBytes())

	return cmp == -1
}

func (a *addedSlice) Swap(i, j int) {
	(*a)[i], (*a)[j] = (*a)[j], (*a)[i]
}

func (a *addedSlice) Get(i int) *row {
	return (*a)[i]
}

func (a *addedSlice) Set(i int, e *row) {
	(*a)[i] = e
}

func (a *addedSlice) SortStableAndRemoveDuplicates() {
	if a.Len() == 0 {
		return
	}

	slices.SortStableFunc(*a, func(x, y *row) int {
		return bytes.Compare(x.left.GetShaBytes(), y.left.GetShaBytes())
	})

	e := rowEqualerShaOnly{}
	last := a.Get(a.Len() - 1)

	for i := a.Len() - 2; i >= 0; i-- {
		x := a.Get(i)

		if e.Equals(last, x) {
			a.Set(i, nil)
		} else {
			last = x
		}
	}

	*a = slices.DeleteFunc(*a, func(n *row) bool {
		return n == nil
	})
}
