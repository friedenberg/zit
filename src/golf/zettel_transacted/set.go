package zettel_transacted

import (
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/hinweis"
)

func MakeSetUnique(c int) Set {
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
		innerMap: make(map[string]Zettel, c),
	}
}

func MakeSetHinweis(c int) Set {
	return Set{
		keyFunc: func(sz Zettel) string {
			return makeKey(sz.Named.Hinweis)
		},
		innerMap: make(map[string]Zettel, c),
	}
}

func (m Set) AddFrom(ch <-chan Zettel) {
	for z := range ch {
		m.Add(z)
	}
}

func (m Set) Get(
	s fmt.Stringer,
) (tz Zettel, ok bool) {
	tz, ok = m.innerMap[s.String()]
	return
}

func (m Set) WriteZettelTransacted(z Zettel) (err error) {
	m.Add(z)

	return
}

func (m Set) Add(z Zettel) {
	k := m.keyFunc(z)

	if _, ok := m.innerMap[k]; ok {
		// errors.Printf("replacing %s with %s", old, z)
	} else {
		// errors.Printf("adding %s", z)
	}

	m.innerMap[k] = z
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

func (a Set) Filter(keyFunc SetKeyFunc, f func(Zettel) (bool, error)) (b Set, err error) {
	if keyFunc == nil {
		keyFunc = a.keyFunc
	}

	b = Set{
		keyFunc:  keyFunc,
		innerMap: make(map[string]Zettel, a.Len()),
	}

	for _, sz := range a.innerMap {
		var ok bool

		ok, err = f(sz)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		if ok {
			b.Add(sz)
		}
	}

	return
}

func (m Set) ToSlice() (s Slice) {
	s = MakeSlice(m.Len())

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

func (s Set) ToSliceHinweisen() (b []hinweis.Hinweis) {
	b = make([]hinweis.Hinweis, 0, s.Len())

	for _, z := range s.innerMap {
		b = append(b, z.Named.Hinweis)
	}

	return
}
