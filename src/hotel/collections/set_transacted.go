package collections

import (
	"fmt"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
)

type SetTransacted struct {
	keyFunc  func(zettel_transacted.Zettel) string
	innerMap map[string]zettel_transacted.Zettel
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

func MakeSetUniqueTransacted(c int) SetTransacted {
	return SetTransacted{
		keyFunc: func(sz zettel_transacted.Zettel) string {
			return makeKey(
				sz.Kopf,
				sz.Mutter,
				sz.Schwanz,
				sz.Named.Hinweis,
				sz.Named.Stored.Sha,
			)
		},
		innerMap: make(map[string]zettel_transacted.Zettel),
	}
}

func MakeSetHinweisTransacted() SetTransacted {
	return SetTransacted{
		keyFunc: func(sz zettel_transacted.Zettel) string {
			return makeKey(sz.Named.Hinweis)
		},
		innerMap: make(map[string]zettel_transacted.Zettel),
	}
}

func (m SetTransacted) Get(
	s fmt.Stringer,
) (tz zettel_transacted.Zettel, ok bool) {
	tz, ok = m.innerMap[s.String()]
	return
}

func (m SetTransacted) Add(z zettel_transacted.Zettel) {
	m.innerMap[m.keyFunc(z)] = z
}

func (m SetTransacted) Del(z zettel_transacted.Zettel) {
	delete(m.innerMap, m.keyFunc(z))
}

func (m SetTransacted) Len() int {
	return len(m.innerMap)
}

func (a SetTransacted) Merge(b SetTransacted) {
	for _, z := range b.innerMap {
		a.Add(z)
	}
}

func (a SetTransacted) Contains(z zettel_transacted.Zettel) bool {
	_, ok := a.innerMap[a.keyFunc(z)]
	return ok
}

func (a SetTransacted) Any() (tz zettel_transacted.Zettel) {
	for _, sz := range a.innerMap {
		tz = sz
		break
	}

	return
}

func (a SetTransacted) Each(f func(zettel_transacted.Zettel) error) (err error) {
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

func (m SetTransacted) ToSlice() (s SliceTransacted) {
	s = MakeSliceTransacted()

	for _, sz := range m.innerMap {
		s.Append(sz)
	}

	return
}

func (s SetTransacted) ToSetPrefixTransacted() (b SetPrefixTransacted) {
	b = MakeSetPrefixTransacted(len(s.innerMap))

	for _, z := range s.innerMap {
		b.Add(z)
	}

	return
}
