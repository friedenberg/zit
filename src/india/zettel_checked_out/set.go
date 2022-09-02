package zettel_checked_out

import (
	"fmt"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/bravo/errors"
)

type Set struct {
	keyFunc  func(CheckedOut) string
	innerMap map[string]CheckedOut
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

func MakeSetUniqueCheckedOut(c int) Set {
	return Set{
		keyFunc: func(sz CheckedOut) string {
			return makeKey(
				sz.Internal.Kopf,
				sz.Internal.Mutter,
				sz.Internal.Schwanz,
				sz.Internal.Named.Hinweis,
				sz.Internal.Named.Stored.Sha,
			)
		},
		innerMap: make(map[string]CheckedOut),
	}
}

func MakeSetHinweisCheckedOut() Set {
	return Set{
		keyFunc: func(sz CheckedOut) string {
			return makeKey(sz.Internal.Named.Hinweis)
		},
		innerMap: make(map[string]CheckedOut),
	}
}

func (m Set) Get(
	s fmt.Stringer,
) (tz CheckedOut, ok bool) {
	tz, ok = m.innerMap[s.String()]
	return
}

func (m Set) Add(z CheckedOut) {
	m.innerMap[m.keyFunc(z)] = z
}

func (m Set) Del(z CheckedOut) {
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

func (a Set) Contains(z CheckedOut) bool {
	_, ok := a.innerMap[a.keyFunc(z)]
	return ok
}

func (a Set) Any() (tz CheckedOut) {
	for _, sz := range a.innerMap {
		tz = sz
		break
	}

	return
}

func (a Set) Each(f func(CheckedOut) error) (err error) {
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

// type MapToKeyFunc struct {
//   Results []string
//   KeyFunc func(CheckedOut) string
// }

// func MakeMapToKeyFunc(f func(CheckedOut) string) MapToKeyFunc {
//   return MapToKeyFunc{
//     Results: make([]string, 0),
//     KeyFunc: f,
//   }
// }

