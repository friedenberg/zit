package collections_ptr

import (
	"sync"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
)

type mutableSetExperimental[
	T schnittstellen.Element,
	TPtr schnittstellen.ElementPtr[T],
] struct {
	K schnittstellen.StringKeyerPtr[T, TPtr]
	E map[string]TPtr

	lAdded          sync.Mutex
	addedDuringIter []TPtr

	lDeleted          sync.Mutex
	deletedDuringIter []string

	l sync.RWMutex
}

func (s *mutableSetExperimental[T, TPtr]) Len() int {
	defer s.TryProcessMutationsDuringReads()
	s.l.RLock()
	defer s.l.RUnlock()

	return s.len()
}

func (s *mutableSetExperimental[T, TPtr]) len() int {
	if s.E == nil {
		return 0
	}

	return len(s.E)
}

func (a *mutableSetExperimental[T, TPtr]) EqualsSetPtrLike(
	b schnittstellen.SetPtrLike[T, TPtr],
) bool {
	defer a.TryProcessMutationsDuringReads()
	a.l.RLock()
	defer a.l.RUnlock()

	return a.EqualsSetLike(b)
}

func (a *mutableSetExperimental[T, TPtr]) EqualsSetLike(
	b schnittstellen.SetLike[T],
) bool {
	if b == nil {
		return false
	}

	defer a.TryProcessMutationsDuringReads()
	a.l.RLock()
	defer a.l.RUnlock()

	if a.len() != b.Len() {
		return false
	}

	for k, va := range a.E {
		vb, ok := b.Get(k)

		if !ok || !va.EqualsAny(vb) {
			return false
		}
	}

	return true
}

func (s *mutableSetExperimental[T, TPtr]) Key(e T) string {
	defer s.TryProcessMutationsDuringReads()
	s.l.RLock()
	defer s.l.RUnlock()

	return s.key(e)
}

func (s *mutableSetExperimental[T, TPtr]) key(e T) string {
	return s.K.GetKey(e)
}

func (s *mutableSetExperimental[T, TPtr]) KeyPtr(e TPtr) string {
	defer s.TryProcessMutationsDuringReads()
	s.l.RLock()
	defer s.l.RUnlock()

	return s.keyPtr(e)
}

func (s *mutableSetExperimental[T, TPtr]) keyPtr(e TPtr) string {
	return s.K.GetKeyPtr(e)
}

func (s *mutableSetExperimental[T, TPtr]) GetPtr(k string) (e TPtr, ok bool) {
	defer s.TryProcessMutationsDuringReads()
	s.l.RLock()
	defer s.l.RUnlock()

	e, ok = s.E[k]

	return
}

func (s *mutableSetExperimental[T, TPtr]) Get(k string) (e T, ok bool) {
	defer s.TryProcessMutationsDuringReads()
	s.l.RLock()
	defer s.l.RUnlock()

	var e1 TPtr

	if e1, ok = s.E[k]; ok {
		e = *e1
	}

	return
}

func (s *mutableSetExperimental[T, TPtr]) ContainsKey(k string) (ok bool) {
	defer s.TryProcessMutationsDuringReads()
	s.l.RLock()
	defer s.l.RUnlock()

	if k == "" {
		return
	}

	_, ok = s.E[k]

	return
}

func (s *mutableSetExperimental[T, TPtr]) Contains(e T) (ok bool) {
	defer s.TryProcessMutationsDuringReads()
	s.l.RLock()
	defer s.l.RUnlock()

	return s.ContainsKey(s.Key(e))
}

func (s *mutableSetExperimental[T, TPtr]) Any() (v T) {
	defer s.TryProcessMutationsDuringReads()
	s.l.RLock()
	defer s.l.RUnlock()

	for _, v1 := range s.E {
		v = *v1
		break
	}

	return
}

// If a read is taking place, this method will block until that read is done.
func (s *mutableSetExperimental[T, TPtr]) ProcessMutationsDuringReads() (err error) {
	s.l.Lock()
	defer s.l.Unlock()

	return s.processMutationsDuringReads()
}

// If a read is taking place, this method will not block until that read is
// done.
func (s *mutableSetExperimental[T, TPtr]) TryProcessMutationsDuringReads() (err error) {
	if !s.l.TryLock() {
		return
	}

	defer s.l.Unlock()

	return s.processMutationsDuringReads()
}

func (s *mutableSetExperimental[T, TPtr]) processMutationsDuringReads() (err error) {
	if !s.l.TryLock() {
		return
	}

	defer s.l.Unlock()

	for i := range s.addedDuringIter {
		if err = s.addPtr(s.addedDuringIter[i]); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	for _, k := range s.deletedDuringIter {
		if err = s.delKey(k); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *mutableSetExperimental[T, TPtr]) Del(v T) (err error) {
	s.l.Lock()
	defer s.l.Unlock()

	return s.delKey(s.key(v))
}

func (s *mutableSetExperimental[T, TPtr]) DelPtr(v TPtr) (err error) {
	k := s.K.GetKeyPtr(v)

	if !s.l.TryLock() {
		s.lDeleted.Lock()
		defer s.lDeleted.Unlock()

		s.deletedDuringIter = append(s.deletedDuringIter, k)

		return
	}

	s.l.Lock()
	defer s.l.Unlock()

	return s.delKey(k)
}

func (s *mutableSetExperimental[T, TPtr]) DelKey(k string) (err error) {
	if !s.l.TryLock() {
		s.lDeleted.Lock()
		defer s.lDeleted.Unlock()

		s.deletedDuringIter = append(s.deletedDuringIter, k)

		return
	}

	defer s.l.Unlock()

	return s.delKey(k)
}

func (s *mutableSetExperimental[T, TPtr]) delKey(k string) (err error) {
	delete(s.E, k)

	return
}

func (s *mutableSetExperimental[T, TPtr]) Add(v T) (err error) {
	if !s.l.TryLock() {
		s.lAdded.Lock()
		defer s.lAdded.Unlock()

		s.addedDuringIter = append(s.addedDuringIter, &v)

		return
	}

	defer s.l.Unlock()

	s.E[s.key(v)] = TPtr(&v)

	return
}

func (s *mutableSetExperimental[T, TPtr]) AddPtr(v TPtr) (err error) {
	if !s.l.TryLock() {
		s.lAdded.Lock()
		defer s.lAdded.Unlock()

		s.addedDuringIter = append(s.addedDuringIter, v)

		return
	}

	defer s.l.Unlock()

	return s.addPtr(v)
}

func (s *mutableSetExperimental[T, TPtr]) addPtr(v TPtr) (err error) {
	s.E[s.K.GetKeyPtr(v)] = v

	return
}

func (s *mutableSetExperimental[T, TPtr]) Elements() (out []T) {
	defer s.TryProcessMutationsDuringReads()
	s.l.RLock()
	defer s.l.RUnlock()

	out = make([]T, 0, s.Len())

	for _, v := range s.E {
		out = append(out, *v)
	}

	return
}

func (s *mutableSetExperimental[T, TPtr]) EachKey(
	wf schnittstellen.FuncIterKey,
) (err error) {
	defer s.TryProcessMutationsDuringReads()
	s.l.RLock()
	defer s.l.RUnlock()

	for v := range s.E {
		if err = wf(v); err != nil {
			if iter.IsStopIteration(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (s *mutableSetExperimental[T, TPtr]) Each(
	wf schnittstellen.FuncIter[T],
) (err error) {
	defer s.TryProcessMutationsDuringReads()
	s.l.RLock()
	defer s.l.RUnlock()

	for _, v := range s.E {
		if err = wf(*v); err != nil {
			if iter.IsStopIteration(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (s *mutableSetExperimental[T, TPtr]) EachPtr(
	wf schnittstellen.FuncIter[TPtr],
) (err error) {
	defer s.TryProcessMutationsDuringReads()
	s.l.RLock()
	defer s.l.RUnlock()

	return s.eachPtr(wf)
}

func (s *mutableSetExperimental[T, TPtr]) eachPtr(
	wf schnittstellen.FuncIter[TPtr],
) (err error) {
	for _, v := range s.E {
		if err = wf(v); err != nil {
			if iter.IsStopIteration(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (a *mutableSetExperimental[T, TPtr]) Reset() {
	if !a.l.TryLock() {
		panic("attempting to reset mutable set during read")
	}

	defer a.l.Unlock()

	for k := range a.E {
		delete(a.E, k)
	}

	a.addedDuringIter = nil
	a.deletedDuringIter = nil
}
