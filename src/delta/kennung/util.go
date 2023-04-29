package kennung

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
)

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
