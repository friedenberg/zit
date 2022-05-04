package etikett

import "github.com/friedenberg/zit/alfa/errors"

type Slice []Etikett

func NewSlice(es ...string) (s Slice, err error) {
	s = make([]Etikett, len(es))

	for i, e := range es {
		if err = s[i].Set(e); err != nil {
			err = errors.Error(err)
			return
		}
	}

	return
}

func (s Slice) ToSet() Set {
	return NewSet([]Etikett(s)...)
}
