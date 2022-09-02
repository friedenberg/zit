package zettel_transacted

import (
	"fmt"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/bravo/errors"
)

type Set struct {
	keyFunc  func(Zettel) string
	innerMap map[string]Zettel
}

func makeKey(ss ...fmt.Stringer) string {
	sb := &strings.Builder{}

	for i, s := range ss {
		if i > 0 {
			sb.WriteString(".")
		}

		sb.WriteString(s.String())
	}

	return sb.String()
}

func MakeSetUniqueTransacted(c int) Set {
	return Set{
		keyFunc: func(sz Zettel) string {
			return makeKey(
				sz.Kopf,
				sz.Mutter,
				sz.Schwanz,
				sz.Named.Hinweis,
				sz.Named.Stored.Sha,
			)
		},
		innerMap: make(map[string]Zettel),
	}
}

func MakeSetHinweisTransacted() Set {
	return Set{
		keyFunc: func(sz Zettel) string {
			return makeKey(sz.Named.Hinweis)
		},
		innerMap: make(map[string]Zettel),
	}
}

func (m Set) Get(
	s fmt.Stringer,
) (tz Zettel, ok bool) {
	tz, ok = m.innerMap[s.String()]
	return
}

func (m Set) Add(z Zettel) {
	m.innerMap[m.keyFunc(z)] = z
}

func (m Set) Del(z Zettel) {
	delete(m.innerMap, m.keyFunc(z))
}

func (m Set) Len() int {
	return len(m.innerMap)
}

func (a Set) Merge(b Set) {
	for _, z := range b.innerMap {
		a.Add(z)
	}
}

func (a Set) Contains(z Zettel) bool {
	_, ok := a.innerMap[a.keyFunc(z)]
	return ok
}

func (a Set) Any() (tz Zettel) {
	for _, sz := range a.innerMap {
		tz = sz
		break
	}

	return
}

func (a Set) Each(f func(Zettel) error) (err error) {
	for _, sz := range a.innerMap {
		if err = f(sz); err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
			} else {
				err = errors.Error(err)
			}

			return
		}
	}

	return
}

func (m Set) ToSlice() (s SliceTransacted) {
	s = MakeSliceTransacted()

	for _, sz := range m.innerMap {
		s.Append(sz)
	}

	return
}

func (s Set) ToSetPrefixTransacted() (b SetPrefixTransacted) {
	b = MakeSetPrefixTransacted(len(s.innerMap))

	for _, z := range s.innerMap {
		b.Add(z)
	}

	return
}
