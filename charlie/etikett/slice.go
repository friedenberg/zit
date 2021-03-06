package etikett

import (
	"sort"
	"strings"

	"github.com/friedenberg/zit/alfa/errors"
)

type Slice []Etikett

func NewSliceFromStrings(es ...string) (s Slice, err error) {
	s = make([]Etikett, len(es))

	for i, e := range es {
		if err = s[i].Set(e); err != nil {
			err = errors.Error(err)
			return
		}
	}

	return
}

func (s Slice) Len() int {
	return len(s)
}

func NewSlice(es ...Etikett) (s Slice) {
	s = make([]Etikett, len(es))

	for i, e := range es {
		s[i] = e
	}

	return
}

func (es *Slice) AddString(v string) (err error) {
	var e Etikett

	if err = e.Set(v); err != nil {
		err = _Error(err)
		return
	}

	es.Add(e)

	return
}

func (es *Slice) Add(e Etikett) {
	*es = append(*es, e)
}

func (s *Slice) Set(v string) (err error) {
	es := strings.Split(v, ",")

	for _, e := range es {
		if err = s.AddString(e); err != nil {
			err = _Error(err)
			return
		}
	}

	return
}

func (es Slice) SortedString() (out []string) {
	out = make([]string, len(es))

	i := 0

	for _, e := range es {
		out[i] = e.String()
		i++
	}

	sort.Slice(
		out,
		func(i, j int) bool {
			return out[i] < out[j]
		},
	)

	return
}

func (s Slice) String() string {
	return strings.Join(s.SortedString(), ", ")
}

func (s Slice) ToSet() Set {
	return MakeSet([]Etikett(s)...)
}
