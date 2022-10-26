package zettel_checked_out

import (
	"fmt"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/golf/zettel_external"
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

func MakeSetUnique(c int) Set {
	return Set{
		keyFunc: func(sz Zettel) string {
			return makeKey(
				sz.Internal.Kopf,
				sz.Internal.Mutter,
				sz.Internal.Schwanz,
				sz.Internal.Named.Hinweis,
				sz.Internal.Named.Stored.Sha,
			)
		},
		innerMap: make(map[string]Zettel, c),
	}
}

func MakeSetHinweisZettel(c int) Set {
	return Set{
		keyFunc: func(sz Zettel) string {
			return makeKey(sz.Internal.Named.Hinweis)
		},
		innerMap: make(map[string]Zettel, c),
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
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (s Set) ToSlice() (out []Zettel) {
	out = make([]Zettel, 0, s.Len())

	for _, z := range s.innerMap {
		out = append(out, z)
	}

	return
}

func (s Set) ToSliceZettelsExternal() (out []zettel_external.Zettel) {
	out = make([]zettel_external.Zettel, 0, s.Len())

	for _, z := range s.innerMap {
		out = append(out, z.External)
	}

	return
}

func (s Set) ToSliceFilesZettelen() (out []string) {
	out = make([]string, 0, s.Len())

	for _, z := range s.innerMap {
		p := z.External.ZettelFD.Path

		if p != "" {
			out = append(out, p)
		}
	}

	return
}

func (s Set) ToSliceFilesAkten() (out []string) {
	out = make([]string, 0, s.Len())

	for _, z := range s.innerMap {
		p := z.External.AkteFD.Path

		if p != "" {
			out = append(out, p)
		}
	}

	return
}
