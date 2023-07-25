package collections

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type MutableSetPtr[T schnittstellen.ValueLike, TPtr schnittstellen.ValuePtr[T]] map[string]TPtr

func MakeMutableSetPtr[T schnittstellen.ValueLike, TPtr schnittstellen.ValuePtr[T]](
	es ...TPtr,
) (s MutableSetPtr[T, TPtr]) {
	s = MutableSetPtr[T, TPtr](make(map[string]TPtr, len(es)))

	for i := range es {
		e := es[i]
		s[e.String()] = e
	}

	return
}

func (s MutableSetPtr[T, TPtr]) Len() int {
	if s == nil {
		return 0
	}

	return len(s)
}

func (a MutableSetPtr[T, TPtr]) EqualsSetLike(
	b schnittstellen.SetLike[T],
) bool {
	if b == nil {
		return false
	}

	if a.Len() != b.Len() {
		return false
	}

	for k, va := range a {
		vb, ok := b.Get(k)

		if !ok || !va.EqualsAny(vb) {
			return false
		}
	}

	return true
}

func (s MutableSetPtr[T, TPtr]) Key(e T) string {
	return e.String()
}

func (s MutableSetPtr[T, TPtr]) Get(k string) (e T, ok bool) {
	var e1 TPtr
	e1, ok = s[k]
	e = *e1
	return
}

func (s MutableSetPtr[T, TPtr]) ContainsKey(k string) (ok bool) {
	if k == "" {
		return
	}

	_, ok = s[k]

	return
}

func (s MutableSetPtr[T, TPtr]) Contains(e T) (ok bool) {
	return s.ContainsKey(s.Key(e))
}

func (s MutableSetPtr[T, TPtr]) Any() (v T) {
	for _, v1 := range s {
		v = *v1
		break
	}

	return
}

func (s MutableSetPtr[T, TPtr]) Del(v T) (err error) {
	return s.DelKey(v.String())
}

func (s MutableSetPtr[T, TPtr]) DelPtr(v TPtr) (err error) {
	return s.DelKey(v.String())
}

func (s MutableSetPtr[T, TPtr]) DelKey(k string) (err error) {
	delete(s, k)
	return
}

func (s MutableSetPtr[T, TPtr]) Add(v T) (err error) {
	s[v.String()] = TPtr(&v)
	return
}

func (s MutableSetPtr[T, TPtr]) AddPtr(v TPtr) (err error) {
	s[v.String()] = v
	return
}

func (s MutableSetPtr[T, TPtr]) Elements() (out []T) {
	out = make([]T, 0, s.Len())

	for _, v := range s {
		out = append(out, *v)
	}

	return
}

func (s MutableSetPtr[T, TPtr]) EachKey(
	wf schnittstellen.FuncIterKey,
) (err error) {
	for v := range s {
		if err = wf(v); err != nil {
			if errors.Is(err, MakeErrStopIteration()) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (s MutableSetPtr[T, TPtr]) Each(
	wf schnittstellen.FuncIter[T],
) (err error) {
	for _, v := range s {
		if err = wf(*v); err != nil {
			if errors.Is(err, MakeErrStopIteration()) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (s MutableSetPtr[T, TPtr]) EachPtr(
	wf schnittstellen.FuncIter[TPtr],
) (err error) {
	for _, v := range s {
		if err = wf(v); err != nil {
			if errors.Is(err, MakeErrStopIteration()) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (a MutableSetPtr[T, TPtr]) Reset() {
	for k := range a {
		delete(a, k)
	}
}

func (a MutableSetPtr[T, TPtr]) CloneSetLike() schnittstellen.SetLike[T] {
	b := MakeSetPtr[T, TPtr]()

	for k, v := range a {
		b[k] = v
	}

	return b
}

func (a MutableSetPtr[T, TPtr]) CloneMutableSetLike() schnittstellen.MutableSetLike[T] {
	c := MakeMutableSetPtr[T, TPtr]()
	a.Each(c.Add)
	return c
}
