package collections

// type set[T schnittstellen.ValueLike] struct {
// 	keyFunc    func(T) string
// 	closed     bool
// 	elementMap map[string]T
// }

// func makeSet[T schnittstellen.ValueLike](kf KeyFunc[T], es ...T) (s *set[T])
// {
// 	t := *new(T)
// 	// Required because interface types do not properly get handled by
// 	// `reflect.TypeOf`
// 	t1 := make([]T, 1)

// 	if reflect.TypeOf(t1).Elem().Kind() == reflect.Interface {
// 		kf(t1[0])
// 	} else {
// 		// confirms that the key function supports nil pointers properly
// 		switch reflect.TypeOf(t).Kind() {
// 		// case reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
// 		case reflect.Ptr:
// 			kf(t)
// 		}
// 	}

// 	s = &set[T]{
// 		keyFunc:    kf,
// 		elementMap: make(map[string]T, len(es)),
// 	}

// 	s.open()
// 	defer s.close()

// 	for _, e := range es {
// 		s.add(e)
// 	}

// 	return s
// }

// func (s *set[T]) open() {
// 	s.closed = false
// }

// func (s *set[T]) close() {
// 	s.closed = true
// }

// func (s set[T]) Len() int {
// 	if s.elementMap == nil {
// 		return 0
// 	}

// 	return len(s.elementMap)
// }

// func (a set[T]) EqualsSetLike(b schnittstellen.SetLike[T]) bool {
// 	if b == nil {
// 		return false
// 	}

// 	if a.Len() != b.Len() {
// 		return false
// 	}

// 	for k, va := range a.elementMap {
// 		vb, ok := b.Get(k)

// 		if !ok || !va.EqualsAny(vb) {
// 			return false
// 		}
// 	}

// 	return true
// }

// func (s set[T]) Key(e T) string {
// 	if s.keyFunc == nil {
// 		return e.String()
// 	} else {
// 		return s.keyFunc(e)
// 	}
// }

// func (s set[T]) Get(k string) (e T, ok bool) {
// 	e, ok = s.elementMap[k]
// 	return
// }

// func (s set[T]) Any() (e T) {
// 	for _, e1 := range s.elementMap {
// 		return e1
// 	}

// 	return
// }

// func (s set[T]) ContainsKey(k string) (ok bool) {
// 	if k == "" {
// 		return
// 	}

// 	_, ok = s.elementMap[k]

// 	return
// }

// func (s set[T]) Contains(e T) (ok bool) {
// 	return s.ContainsKey(s.Key(e))
// }

// func (es set[T]) add(e T) (err error) {
// 	if es.closed {
// 		panic(fmt.Sprintf("trying to add %T to closed set", e))
// 	}

// 	es.elementMap[es.Key(e)] = e

// 	return
// }

// func (es set[T]) addCustom(e T, kf func(T) string) (err error) {
// 	if es.closed {
// 		panic(fmt.Sprintf("trying to add %T to closed set", e))
// 	}

// 	if kf == nil {
// 		kf = es.Key
// 	}

// 	es.elementMap[kf(e)] = e

// 	return
// }

// func (s set[T]) EachKey(wf schnittstellen.FuncIterKey) (err error) {
// 	for v := range s.elementMap {
// 		if err = wf(v); err != nil {
// 			if errors.Is(err, MakeErrStopIteration()) {
// 				err = nil
// 			} else {
// 				err = errors.Wrap(err)
// 			}

// 			return
// 		}
// 	}

// 	return
// }

// func (s set[T]) Elements() (out []T) {
// 	out = make([]T, 0, s.Len())

// 	for _, v := range s.elementMap {
// 		out = append(out, v)
// 	}

// 	return
// }

// func (s set[T]) Each(wf schnittstellen.FuncIter[T]) (err error) {
// 	for _, v := range s.elementMap {
// 		if err = wf(v); err != nil {
// 			if errors.Is(err, MakeErrStopIteration()) {
// 				err = nil
// 			} else {
// 				err = errors.Wrap(err)
// 			}

// 			return
// 		}
// 	}

// 	return
// }

// func (s set[T]) EachPtr(wf schnittstellen.FuncIter[*T]) (err error) {
// 	for _, v := range s.elementMap {
// 		if err = wf(&v); err != nil {
// 			if errors.Is(err, MakeErrStopIteration()) {
// 				err = nil
// 			} else {
// 				err = errors.Wrap(err)
// 			}

// 			return
// 		}
// 	}

// 	return
// }

// func (a set[T]) CloneSetLike() schnittstellen.SetLike[T] {
// 	c := makeSet(a.Key)
// 	c.open()
// 	defer c.close()

// 	a.Each(c.add)

// 	return c
// }

// func (a set[T]) CloneMutableSetLike() schnittstellen.MutableSetLike[T] {
// 	c := MakeMutableSet(a.Key)
// 	a.Each(c.Add)
// 	return c
// }

// func (s set[T]) MarshalBinary() (bs []byte, err error) {
// 	b := bytes.NewBuffer(bs)
// 	enc := gob.NewEncoder(b)

// 	if err = enc.Encode(s.elementMap); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	bs = b.Bytes()

// 	return
// }

// func (s *set[T]) UnmarshalBinary(bs []byte) (err error) {
// 	b := bytes.NewBuffer(bs)
// 	dec := gob.NewDecoder(b)

// 	if err = dec.Decode(&s.elementMap); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }
