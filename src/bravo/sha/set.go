package sha

import (
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type Set struct {
	innerMap map[string]Sha
}

func MakeSet(c int) Set {
	return Set{
		innerMap: make(map[string]Sha),
	}
}

func (m Set) Get(
	s fmt.Stringer,
) (sh Sha, ok bool) {
	sh, ok = m.innerMap[s.String()]
	return
}

func (m Set) Add(sh Sha) {
	m.innerMap[sh.String()] = sh
}

func (m Set) Del(sh Sha) {
	delete(m.innerMap, sh.String())
}

func (m Set) Len() int {
	return len(m.innerMap)
}

func (a Set) Merge(b Set) {
	for _, z := range b.innerMap {
		a.Add(z)
	}
}

func (a Set) Contains(sh Sha) bool {
	_, ok := a.innerMap[sh.String()]
	return ok
}

func (a Set) Any() (sh Sha) {
	for _, sh = range a.innerMap {
		break
	}

	return
}

func (a Set) Each(f func(Sha) error) (err error) {
	for _, sh := range a.innerMap {
		if err = f(sh); err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (a Set) Filter(f func(Sha) (bool, error)) (b Set, err error) {
	b = Set{
		innerMap: make(map[string]Sha, a.Len()),
	}

	for _, sh := range a.innerMap {
		var ok bool

		ok, err = f(sh)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		if ok {
			b.Add(sh)
		}
	}

	return
}

func (m Set) ToSlice() (s Slice) {
	s = MakeSlice(m.Len())

	for _, sh := range m.innerMap {
		s.Append(sh)
	}

	return
}
