package collections

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type SetPtr[T schnittstellen.ValueLike, TPtr schnittstellen.ValuePtr[T]] map[string]TPtr

func MakeSetPtrValueCustom[T schnittstellen.ValueLike, TPtr schnittstellen.ValuePtr[T]](
	kf func(T) string,
	es ...T,
) (s SetPtr[T, TPtr]) {
	s = SetPtr[T, TPtr](make(map[string]TPtr, len(es)))

	for i := range es {
		e := TPtr(&es[i])
		s[kf(T(*e))] = e
	}

	return
}

func MakeSetPtrCustom[T schnittstellen.ValueLike, TPtr schnittstellen.ValuePtr[T]](
	kf func(T) string,
	es ...TPtr,
) (s SetPtr[T, TPtr]) {
	s = SetPtr[T, TPtr](make(map[string]TPtr, len(es)))

	for i := range es {
		e := es[i]
		s[kf(T(*e))] = e
	}

	return
}

func MakeSetPtrValue[T schnittstellen.ValueLike, TPtr schnittstellen.ValuePtr[T]](
	es ...T,
) (s SetPtr[T, TPtr]) {
	s = SetPtr[T, TPtr](make(map[string]TPtr, len(es)))

	for i := range es {
		e := TPtr(&es[i])
		s[e.String()] = e
	}

	return
}

func MakeSetPtr[T schnittstellen.ValueLike, TPtr schnittstellen.ValuePtr[T]](
	es ...TPtr,
) (s SetPtr[T, TPtr]) {
	s = SetPtr[T, TPtr](make(map[string]TPtr, len(es)))

	for i := range es {
		e := es[i]
		s[e.String()] = e
	}

	return
}

func (s SetPtr[T, TPtr]) Len() int {
	if s == nil {
		return 0
	}

	return len(s)
}

func (a SetPtr[T, TPtr]) EqualsSetLike(b schnittstellen.SetLike[T]) bool {
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

func (s SetPtr[T, TPtr]) Key(e T) string {
	return e.String()
}

func (s SetPtr[T, TPtr]) Get(k string) (e T, ok bool) {
	var e1 TPtr

	if e1, ok = s[k]; ok {
		e = *e1
	}

	return
}

func (s SetPtr[T, TPtr]) Any() (e T) {
	for _, e1 := range s {
		return *e1
	}

	return
}

func (s SetPtr[T, TPtr]) ContainsKey(k string) (ok bool) {
	if k == "" {
		return
	}

	_, ok = s[k]

	return
}

func (s SetPtr[T, TPtr]) Contains(e T) (ok bool) {
	return s.ContainsKey(s.Key(e))
}

func (s SetPtr[T, TPtr]) EachKey(wf schnittstellen.FuncIterKey) (err error) {
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

func (s SetPtr[T, TPtr]) Elements() (out []T) {
	out = make([]T, 0, s.Len())

	for _, v := range s {
		out = append(out, *v)
	}

	return
}

func (s SetPtr[T, TPtr]) Add(v T) (err error) {
	s[v.String()] = TPtr(&v)
	return
}

func (s SetPtr[T, TPtr]) AddPtr(v TPtr) (err error) {
	s[v.String()] = v
	return
}

func (s SetPtr[T, TPtr]) Each(wf schnittstellen.FuncIter[T]) (err error) {
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

func (s SetPtr[T, TPtr]) EachPtr(
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

func (a SetPtr[T, TPtr]) CloneSetLike() schnittstellen.SetLike[T] {
	return a
}

func (a SetPtr[T, TPtr]) CloneMutableSetLike() schnittstellen.MutableSetLike[T] {
	c := MakeMutableSetPtr[T, TPtr]()
	a.Each(c.Add)
	return c
}
