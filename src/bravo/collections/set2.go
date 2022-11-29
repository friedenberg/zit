package collections

import (
	"bytes"
	"encoding/gob"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
)

type Set2[T gattung.Element, T1 gattung.ElementPtr[T]] struct {
	private setPrivate[T, T1]
}

func Set2FromSlice[T gattung.Element, T1 gattung.ElementPtr[T]](
	keyer gattung.Keyer[T, T1],
	es ...T1,
) (s Set2[T, T1]) {
	s.private = setPrivateFromSlice(keyer, es...)

	return
}

func Set2FromSetLike[T gattung.Element, T1 gattung.ElementPtr[T]](
	keyer gattung.Keyer[T, T1],
	s1 SetLike[T1],
) (s Set2[T, T1]) {
	s.private = setPrivateFromSetLike(keyer, s1)

	return
}

func (s Set2[T, T1]) Len() int {
	return len(s.private.Elements)
}

func (s Set2[T, T1]) Get(k string) (e T1, ok bool) {
	e, ok = s.private.Elements[k]
	return
}

func (s Set2[T, T1]) ContainsKey(k string) (ok bool) {
	if k == "" {
		return
	}

	_, ok = s.private.Elements[k]

	return
}

func (s Set2[T, T1]) Contains(e T1) (ok bool) {
	return s.ContainsKey(s.private.Key(e))
}

func (s Set2[T, T1]) EachKey(wf WriterFuncKey) (err error) {
	for v, _ := range s.private.Elements {
		if err = wf(v); err != nil {
			if errors.IsEOF(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (s Set2[T, T1]) Each(wf WriterFunc[T1]) (err error) {
	for _, v := range s.private.Elements {
		if err = wf(v); err != nil {
			if errors.IsEOF(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (s Set2[T, T1]) Elements() (out []T1) {
	out = make([]T1, 0, s.Len())

	for _, v := range s.private.Elements {
		out = append(out, v)
	}

	return
}

func (s *Set2[T, T1]) GobDecode(bs []byte) (err error) {
	b := bytes.NewBuffer(bs)
	dec := gob.NewDecoder(b)

	if err = dec.Decode(&s.private); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Set2[T, T1]) GobEncode() (bs []byte, err error) {
	b := bytes.NewBuffer(bs)
	enc := gob.NewEncoder(b)

	if err = enc.Encode(s.private); err != nil {
		err = errors.Wrap(err)
		return
	}

	bs = b.Bytes()

	return
}
