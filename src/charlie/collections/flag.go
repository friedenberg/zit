package collections

import (
	"flag"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
)

type flagPtr[T schnittstellen.ValueLike] interface {
	schnittstellen.ValuePtr[T]
	schnittstellen.SetterPtr[T]
}

type Flag[
	T schnittstellen.ValueLike,
	TPtr flagPtr[T],
] interface {
	flag.Value
	SetMany(vs ...string) (err error)
	GetSet() schnittstellen.SetLike[T]
	GetMutableSet() schnittstellen.MutableSetLike[T]
}

func MakeFlagCommasFromExisting[
	T schnittstellen.ValueLike,
	TPtr flagPtr[T],
](
	p SetterPolicy,
	existing *schnittstellen.SetLike[T],
) Flag[T, TPtr] {
	if *existing == nil {
		e := schnittstellen.SetLike[T](
			MakeSet(
				(T).String,
			),
		)

		*existing = e
	}

	return &flagCommas[T, TPtr]{
		SetterPolicy: p,
		set:          existing,
	}
}

func MakeFlagCommas[
	T schnittstellen.ValueLike,
	TPtr flagPtr[T],
](p SetterPolicy,
) Flag[T, TPtr] {
	var s schnittstellen.SetLike[T]

	s = MakeMutableSet(
		func(e T) string {
			return e.String()
		},
	)

	return &flagCommas[T, TPtr]{
		SetterPolicy: p,
		set:          &s,
	}
}

type flagCommas[
	T schnittstellen.ValueLike,
	TPtr flagPtr[T],
] struct {
	SetterPolicy
	set *schnittstellen.SetLike[T]
}

func (f flagCommas[T, TPtr]) GetSet() (s schnittstellen.SetLike[T]) {
	return (*f.set).CloneSetLike()
}

func (f flagCommas[T, TPtr]) GetMutableSet() (s schnittstellen.MutableSetLike[T]) {
	return (*f.set).CloneMutableSetLike()
}

func (f flagCommas[T, TPtr]) String() (out string) {
	if f.set == nil {
		return
	}

	sorted := iter.SortedStrings[T](*f.set)

	sb := &strings.Builder{}
	first := true

	for _, e1 := range sorted {
		if !first {
			sb.WriteString(", ")
		}

		sb.WriteString(e1)

		first = false
	}

	out = sb.String()

	return
}

func (f *flagCommas[T, TPtr]) SetMany(vs ...string) (err error) {
	for _, v := range vs {
		if err = f.Set(v); err != nil {
			return
		}
	}

	return
}

func (f *flagCommas[T, TPtr]) Set(v string) (err error) {
	r := (*f.set).CloneMutableSetLike()

	switch f.SetterPolicy {
	case SetterPolicyReset:
		r.Reset()
	}

	els := strings.Split(v, ",")

	for _, e := range els {
		e = strings.TrimSpace(e)

		if err = iter.AddString[T, TPtr](r, e); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	*f.set = r.CloneSetLike()

	return
}
