package collections_ptr

import (
	"flag"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
)

type SetterPolicy int

const (
	SetterPolicyAppend = SetterPolicy(iota)
	SetterPolicyReset
)

type flagPtr[T schnittstellen.ValueLike] interface {
	schnittstellen.ValuePtr[T]
	schnittstellen.SetterPtr[T]
}

// TODO-P2 add Resetter2 and Pool
type Flag[
	T schnittstellen.ValueLike,
	TPtr flagPtr[T],
] interface {
	flag.Value
	SetMany(vs ...string) (err error)
	schnittstellen.MutableSetPtrLike[T, TPtr]
	GetSetPtrLike() schnittstellen.SetPtrLike[T, TPtr]
	GetMutableSetPtrLike() schnittstellen.MutableSetPtrLike[T, TPtr]
}

func MakeFlagCommasFromExisting[
	T schnittstellen.ValueLike,
	TPtr flagPtr[T],
](
	p SetterPolicy,
	existing schnittstellen.MutableSetPtrLike[T, TPtr],
	// pool schnittstellen.Pool[T, TPtr],
	// resetter schnittstellen.Resetter2[T, TPtr],
) Flag[T, TPtr] {
	return &flagCommas[T, TPtr]{
		SP:                p,
		MutableSetPtrLike: existing,
		// pool:              pool,
		// resetter:          resetter,
	}
}

func MakeFlagCommas[
	T schnittstellen.ValueLike,
	TPtr flagPtr[T],
](
	p SetterPolicy,
	// pool schnittstellen.Pool[T, TPtr],
	// resetter schnittstellen.Resetter2[T, TPtr],
) Flag[T, TPtr] {
	return &flagCommas[T, TPtr]{
		SP:                p,
		MutableSetPtrLike: MakeMutableValueSet[T, TPtr](nil),
		// pool:              pool,
		// resetter:          resetter,
	}
}

type flagCommas[
	T schnittstellen.ValueLike,
	TPtr flagPtr[T],
] struct {
	SP SetterPolicy
	schnittstellen.MutableSetPtrLike[T, TPtr]
	pool     schnittstellen.Pool[T, TPtr]
	resetter schnittstellen.Resetter2[T, TPtr]
}

func (f flagCommas[T, TPtr]) GetSetPtrLike() (s schnittstellen.SetPtrLike[T, TPtr]) {
	return f.CloneSetPtrLike()
}

func (f flagCommas[T, TPtr]) GetMutableSetPtrLike() (s schnittstellen.MutableSetPtrLike[T, TPtr]) {
	return f.CloneMutableSetPtrLike()
}

func (f flagCommas[T, TPtr]) String() (out string) {
	if f.MutableSetPtrLike == nil {
		return
	}

	sorted := iter.SortedStrings[T](f)

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
	switch f.SP {
	case SetterPolicyReset:
		f.Reset()
	}

	els := strings.Split(v, ",")

	for _, e := range els {
		e = strings.TrimSpace(e)

		// TODO-P2 use iter.AddStringPtr
		if err = iter.AddString[T, TPtr](f, e); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
