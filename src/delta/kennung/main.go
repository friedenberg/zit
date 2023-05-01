package kennung

import (
	"encoding"
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
)

type QueryPrefixer interface {
	GetQueryPrefix() string
}

type IdLike interface {
	schnittstellen.IdLike
	encoding.TextMarshaler
	encoding.BinaryMarshaler
	Parts() [3]string
}

type KennungPtr interface {
	IdLike
	encoding.TextUnmarshaler
	encoding.BinaryUnmarshaler
	schnittstellen.Resetter2
}

type KennungLike[T any] interface {
	IdLike
	schnittstellen.GattungGetter
	schnittstellen.ValueLike
	schnittstellen.Equatable[T]
}

type KennungLikePtr[T schnittstellen.Value[T]] interface {
	KennungLike[T]
	KennungPtr
	schnittstellen.ValuePtrLike
	schnittstellen.SetterPtr[T]
	schnittstellen.Resetable[T]
}

func Make(v string) (k IdLike, err error) {
	{
		var e Etikett

		if err = e.Set(v); err == nil {
			k = e
			return
		}
	}

	{
		var t Typ

		if err = t.Set(v); err == nil {
			k = t
			return
		}
	}

	{
		var ka Kasten

		if err = ka.Set(v); err == nil {
			k = ka
			return
		}
	}

	{
		var h Hinweis

		if err = h.Set(v); err == nil {
			k = h
			return
		}
	}

	err = errors.Errorf("%q is not a valid Kennung", v)

	return
}

func AlignedParts(
	id IdLike,
	lenLeft, lenRight int,
) (string, string, string) {
	parts := id.Parts()
	left := parts[0]
	middle := parts[1]
	right := parts[2]

	diffLeft := lenLeft - len(left)
	if diffLeft > 0 {
		left = strings.Repeat(" ", diffLeft) + left
	}

	diffRight := lenRight - len(right)
	if diffRight > 0 {
		right = right + strings.Repeat(" ", diffRight)
	}

	return left, middle, right
}

func Aligned(id IdLike, lenLeft, lenRight int) string {
	left, middle, right := AlignedParts(id, lenLeft, lenRight)
	return fmt.Sprintf("%s%s%s", left, middle, right)
}

func LeftSubtract[T schnittstellen.Stringer, TPtr schnittstellen.StringSetterPtr[T]](
	a, b T,
) (c T, err error) {
	if err = TPtr(&c).Set(strings.TrimPrefix(a.String(), b.String())); err != nil {
		err = errors.Wrapf(err, "'%s' - '%s'", a, b)
		return
	}

	return
}

func Contains[T schnittstellen.Stringer](a, b T) bool {
	if len(b.String()) > len(a.String()) {
		return false
	}

	return strings.HasPrefix(a.String(), b.String())
}

func Includes[T schnittstellen.Stringer](a, b T) bool {
	return Contains(b, a)
}

func Less[T schnittstellen.Stringer](a, b T) bool {
	return a.String() < b.String()
}

func LessLen[T schnittstellen.Stringer](a, b T) bool {
	return len(a.String()) < len(b.String())
}

func IsEmpty[T schnittstellen.Stringer](a T) bool {
	return len(a.String()) == 0
}

func SansPrefix(a Etikett) (b Etikett) {
	b = MustEtikett(strings.TrimPrefix(a.String(), "-"))
	return
}

func IsDependentLeaf(a Etikett) (has bool) {
	has = strings.HasPrefix(strings.TrimSpace(a.String()), "-")
	return
}

func HasParentPrefix(a, b Etikett) (has bool) {
	has = strings.HasPrefix(strings.TrimSpace(a.String()), b.String())
	return
}

func IntersectPrefixes(s1 EtikettSet, s2 EtikettSet) (s3 EtikettSet) {
	s4 := MakeEtikettMutableSet()

	for _, e1 := range s2.Elements() {
		didAdd := false

		for _, e := range s1.Elements() {
			if strings.HasPrefix(e.String(), e1.String()) {
				didAdd = true
				s4.Add(e)
			}
		}

		if !didAdd {
			s3 = MakeEtikettMutableSet()
			return
		}
	}

	s3 = s4.ImmutableClone()

	return
}

func SubtractPrefix(s1 EtikettSet, e Etikett) (s2 EtikettSet) {
	s3 := MakeEtikettMutableSet()

	for _, e1 := range s1.Elements() {
		e2, _ := LeftSubtract(e1, e)

		if e2.String() == "" {
			continue
		}

		s3.Add(e2)
	}

	s2 = s3.ImmutableClone()

	return
}

func Description(s EtikettSet) string {
	return collections.StringCommaSeparated[Etikett](s)
}

func WithRemovedCommonPrefixes(s EtikettSet) (s2 EtikettSet) {
	es1 := collections.SortedValues(s)
	es := make([]Etikett, 0, len(es1))

	for _, e := range es1 {
		if len(es) == 0 {
			es = append(es, e)
			continue
		}

		idxLast := len(es) - 1
		last := es[idxLast]

		switch {
		case Contains(last, e):
			continue

		case Contains(e, last):
			es[idxLast] = e

		default:
			es = append(es, e)
		}
	}

	s2 = MakeEtikettSet(es...)

	return
}

func expandOne[T KennungLike[T], TPtr KennungLikePtr[T]](
	k T,
	ex Expander,
	acc schnittstellen.Adder[T],
) {
	f := collections.MakeFuncSetString[T, TPtr](acc)
	ex.Expand(f, k.String())
}

func ExpandOneSlice[T KennungLike[T], TPtr KennungLikePtr[T]](
	k T,
	exes ...Expander,
) (out []T) {
	s1 := collections.MakeMutableSetStringer[T]()

	if len(exes) == 0 {
		exes = []Expander{ExpanderAll}
	}

	for _, ex := range exes {
		expandOne[T, TPtr](k, ex, s1)
	}

	out = collections.SortedValuesBy[T](
		s1,
		func(a, b T) bool {
			return len(a.String()) < len(b.String())
		},
	)

	return
}

func ExpandOne[T KennungLike[T], TPtr KennungLikePtr[T]](
	k T,
	exes ...Expander,
) (out schnittstellen.Set[T]) {
	s1 := collections.MakeMutableSetStringer[T]()

	if len(exes) == 0 {
		exes = []Expander{ExpanderAll}
	}

	for _, ex := range exes {
		expandOne[T, TPtr](k, ex, s1)
	}

	out = s1.ImmutableClone()

	return
}

func ExpandMany[T KennungLike[T], TPtr KennungLikePtr[T]](
	ks schnittstellen.Set[T],
	exes ...Expander,
) (out schnittstellen.Set[T]) {
	s1 := collections.MakeMutableSetStringer[T]()

	if len(exes) == 0 {
		exes = []Expander{ExpanderAll}
	}

	ks.Each(
		func(k T) (err error) {
			for _, ex := range exes {
				expandOne[T, TPtr](k, ex, s1)
			}

			return
		},
	)

	out = s1.ImmutableClone()

	return
}

func Expanded(s EtikettSet, exes ...Expander) (out EtikettSet) {
	return ExpandMany[Etikett, *Etikett](s, exes...)
}

func AddNormalized(es EtikettMutableSet, e Etikett) {
	ExpandOne(e, ExpanderRight).Each(es.Add)
	es.Add(e)

	c := es.ImmutableClone()
	es.Reset()
	WithRemovedCommonPrefixes(c).Each(es.Add)
}

func RemovePrefixes(es EtikettMutableSet, needle Etikett) {
	for _, haystack := range es.Elements() {
		// TODO-P2 make more efficient
		if strings.HasPrefix(haystack.String(), needle.String()) {
			es.Del(haystack)
		}
	}
}

func Withdraw(s1 EtikettMutableSet, e Etikett) (s2 EtikettSet) {
	s3 := MakeEtikettMutableSet()

	for _, e1 := range s1.Elements() {
		if Contains(e1, e) {
			s3.Add(e1)
		}
	}

	s3.Each(s1.Del)
	s2 = s3.ImmutableClone()

	return
}

func MatchTwoSortedEtikettStringSlices(a, b []string) (hasMatch bool) {
	var longer, shorter []string

	switch {
	case len(a) < len(b):
		shorter = a
		longer = b

	default:
		shorter = b
		longer = a
	}

	for _, v := range shorter {
		c := rune(v[0])

		var idx int
		idx, hasMatch = BinarySearchForRuneInEtikettenSortedStringSlice(longer, c)
		errors.Err().Print(idx, hasMatch)

		switch {
		case hasMatch:
			return

		case idx > len(longer)-1:
			return
		}

		longer = longer[idx:]
	}

	return
}

func BinarySearchForRuneInEtikettenSortedStringSlice(
	haystack []string,
	needle rune,
) (idx int, ok bool) {
	var low, hi int
	hi = len(haystack) - 1

	for {
		idx = ((hi - low) / 2) + low
		midValRaw := haystack[idx]

		if midValRaw == "" {
			return
		}

		midVal := rune(midValRaw[0])

		if hi == low {
			ok = midVal == needle
			return
		}

		switch {
		case midVal > needle:
			// search left
			hi = idx - 1
			continue

		case midVal == needle:
			// found
			ok = true
			return

		case midVal < needle:
			// search right
			low = idx + 1
			continue
		}
	}
}
