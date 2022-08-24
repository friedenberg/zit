package collections

import (
	"fmt"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
)

type SetTransacted struct {
	keyFunc  func(stored_zettel.Transacted) string
	innerMap map[string]stored_zettel.Transacted
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

func MakeSetUniqueTransacted() SetTransacted {
	return SetTransacted{
		keyFunc: func(sz stored_zettel.Transacted) string {
			return makeKey(
				sz.Kopf,
				sz.Mutter,
				sz.Schwanz,
				sz.Hinweis,
				sz.Stored.Sha,
			)
		},
		innerMap: make(map[string]stored_zettel.Transacted),
	}
}

func MakeSetHinweisTransacted() SetTransacted {
	return SetTransacted{
		keyFunc: func(sz stored_zettel.Transacted) string {
			return makeKey(sz.Hinweis)
		},
		innerMap: make(map[string]stored_zettel.Transacted),
	}
}

func (m SetTransacted) Get(
	s fmt.Stringer,
) (tz stored_zettel.Transacted, ok bool) {
	tz, ok = m.innerMap[s.String()]
	return
}

func (m SetTransacted) Add(z stored_zettel.Transacted) {
	m.innerMap[m.keyFunc(z)] = z
}

func (m SetTransacted) Del(z stored_zettel.Transacted) {
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

func (a SetTransacted) Contains(z stored_zettel.Transacted) bool {
	_, ok := a.innerMap[a.keyFunc(z)]
	return ok
}

func (a SetTransacted) Any() (tz stored_zettel.Transacted) {
	for _, sz := range a.innerMap {
		tz = sz
		break
	}

	return
}

func (a SetTransacted) Each(f func(stored_zettel.Transacted) error) (err error) {
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
	s = make(SliceTransacted, 0, len(m.innerMap))

	for _, sz := range m.innerMap {
		s = append(s, sz)
	}

	return
}
