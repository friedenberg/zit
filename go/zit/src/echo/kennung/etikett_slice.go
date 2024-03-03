package kennung

import (
	"sort"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
)

// TODO-P3 make mutable / immutable?
type EtikettSlice []Etikett

func MakeEtikettSlice(es ...Etikett) (s EtikettSlice) {
	s = make([]Etikett, len(es))

	for i, e := range es {
		s[i] = e
	}

	return
}

func NewSliceFromStrings(es ...string) (s EtikettSlice, err error) {
	s = make([]Etikett, len(es))

	for i, e := range es {
		if err = s[i].Set(e); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s EtikettSlice) Len() int {
	return len(s)
}

func (es *EtikettSlice) AddString(v string) (err error) {
	var e Etikett

	if err = e.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	es.Add(e)

	return
}

func (es *EtikettSlice) Add(e Etikett) {
	*es = append(*es, e)
}

func (s *EtikettSlice) Set(v string) (err error) {
	es := strings.Split(v, ",")

	for _, e := range es {
		if err = s.AddString(e); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (es EtikettSlice) SortedString() (out []string) {
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

func (s EtikettSlice) String() string {
	return strings.Join(s.SortedString(), ", ")
}

func (s EtikettSlice) ToSet() EtikettSet {
	return MakeEtikettSet([]Etikett(s)...)
}
