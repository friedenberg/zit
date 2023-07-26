package collections

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type setPrivate[T any, T1 schnittstellen.Ptr[T]] struct {
	Elements map[string]T1
	schnittstellen.KeyPtrer[T, T1]
}

func setPrivateFromSetLike[T any, T1 schnittstellen.Ptr[T]](
	keyer schnittstellen.KeyPtrer[T, T1],
	s1 schnittstellen.SetLike[T1],
) (s setPrivate[T, T1]) {
	l := 0

	if s1 != nil {
		l = s1.Len()
	}

	s = setPrivate[T, T1]{
		KeyPtrer: keyer,
		Elements: make(map[string]T1, l),
	}

	// confirms that the key function supports nil pointers properly
	s.Key(nil)

	if s1 != nil {
		s1.Each(s.add)
	}

	return
}

func setPrivateFromSlice[T any, T1 schnittstellen.Ptr[T]](
	keyer schnittstellen.KeyPtrer[T, T1],
	es ...T1,
) (s setPrivate[T, T1]) {
	s = setPrivate[T, T1]{
		KeyPtrer: keyer,
		Elements: make(map[string]T1, len(es)),
	}

	// confirms that the key function supports nil pointers properly
	s.Key(nil)

	for _, e := range es {
		s.add(e)
	}

	return
}

func (es setPrivate[T, T1]) add(e T1) (err error) {
	if e == nil {
		panic(ErrNilPointer)
	}

	es.Elements[es.Key(e)] = e

	return
}
