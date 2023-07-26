package kennung

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections"
)

type QueryPrefixer interface {
	GetQueryPrefix() string
}

type KennungSansGattung interface {
	schnittstellen.ValueLike
	Parts() [3]string
	KennungSansGattungClone() KennungSansGattung
	KennungSansGattungPtrClone() KennungSansGattungPtr
}

type Kennung interface {
	KennungSansGattung
	schnittstellen.GattungGetter
	KennungClone() Kennung
	KennungPtrClone() KennungPtr
}

type KennungSansGattungPtr interface {
	KennungSansGattung
	schnittstellen.Resetter
	schnittstellen.Setter
}

type KennungPtr interface {
	Kennung
	KennungSansGattungPtr
	KennungPtrClone() KennungPtr
}

type KennungLike[T any] interface {
	Kennung
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

type IndexedLike[T KennungSansGattung] interface {
	GetInt() int
	GetKennung() T
	GetSchwanzenCount() int
	GetCount() int
	GetTridex() schnittstellen.Tridex
	GetExpandedRight() schnittstellen.Set[T]
	GetExpandedAll() schnittstellen.Set[T]
}

type Index struct {
	Etiketten func(Etikett) (IndexedLike[Etikett], error)
}

func MakeWithGattung(
	gg schnittstellen.GattungGetter,
	v string,
) (k Kennung, err error) {
	switch gattung.Make(gg.GetGattung()) {
	case gattung.Zettel:
		var h Hinweis

		if err = h.Set(v); err == nil {
			k = h
			return
		}

	case gattung.Etikett:
		var e Etikett

		if err = e.Set(v); err == nil {
			k = e
			return
		}

	case gattung.Typ:
		var t Typ

		if err = t.Set(v); err == nil {
			k = t
			return
		}

	case gattung.Kasten:
		var ka Kasten

		if err = ka.Set(v); err == nil {
			k = ka
			return
		}

	case gattung.Konfig:
		var h Konfig

		if err = h.Set(v); err == nil {
			k = h
			return
		}

	case gattung.Bestandsaufnahme:
		var h Tai

		if err = h.Set(v); err == nil {
			k = h
			return
		}
	}

	err = errors.Errorf("%q is not a valid Kennung", v)

	return
}

func Make(v string) (k Kennung, err error) {
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

	{
		var h Konfig

		if err = h.Set(v); err == nil {
			k = h
			return
		}
	}

	{
		var h Tai

		if err = h.Set(v); err == nil {
			k = h
			return
		}
	}

	err = errors.Errorf("%q is not a valid Kennung", v)

	return
}

func MakeMatcher(
	k KennungSansGattungPtr,
	v string,
	expander func(string) (string, error),
	ki Index,
) (m Matcher, isNegated bool, isExact bool, err error) {
	v = strings.TrimSpace(v)
	didExpand := false

	if expander != nil {
		v1 := v

		if v1, err = expander(v); err != nil {
			err = nil
			v1 = v
		} else {
			didExpand = true
		}

		v = v1
	}

	if isNegated, isExact, err = SetQueryKennung(k, v); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !didExpand && expander != nil {
		v1 := k.String()

		if v1, err = expander(v1); err != nil {
			err = nil
			v1 = v
		}

		if err = k.Set(v1); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if isExact {
		m = MakeMatcherContainsExactly(k)
	} else {
		m = MakeMatcherContains(k, ki)
	}

	if isNegated {
		m = MakeMatcherNegate(m)
	}

	return
}

func SetQueryKennung(
	k schnittstellen.Setter,
	v string,
) (isNegated bool, isExact bool, err error) {
	v = strings.TrimSpace(v)

	if len(v) > 0 && []rune(v)[0] == QueryNegationOperator {
		v = v[1:]
		isNegated = true
	}

	if len(v) > 0 && []rune(v)[len(v)-1] == QueryExactOperator {
		v = v[:len(v)-1]
		isExact = true
	}

	var p string

	if qp, ok := k.(QueryPrefixer); ok {
		p = qp.GetQueryPrefix()
	}

	if len(v) > 0 && v[:len(p)] == p {
		v = v[len(p):]
	}

	if err = k.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func KennungContainsExactlyMatchable(k KennungSansGattung, m Matchable) bool {
	switch kt := k.(type) {
	case EtikettLike:
		es := m.GetEtiketten()

		if es.Contains(kt.GetEtikett()) {
			return true
		}

	case TypLike:
		if ContainsExactly(m.GetTyp(), k) {
			return true
		}

	default:
		// nop
	}

	idl := m.GetIdLike()

	if !ContainsExactly(idl, k) {
		return false
	}

	return true
}

func KennungContainsMatchable(
	k KennungSansGattung,
	m Matchable,
	ki Index,
) bool {
	switch kt := k.(type) {
	case EtikettLike:
		if iter.CheckAny[Etikett](
			m.GetEtiketten(),
			func(e Etikett) (ok bool) {
				indexed, err := ki.Etiketten(e)
				var expanded schnittstellen.Set[Etikett]

				if err == nil {
					expanded = indexed.GetExpandedRight()
				} else {
					expanded = ExpandOne(e, ExpanderRight)
				}

				ok = expanded.Contains(kt.GetEtikett())

				return
			},
		) {
			return true
		}

	case TypLike:
		if Contains(m.GetTyp(), k) {
			return true
		}

	case ShaLike:
		if kt.ContainsMatchable(m) {
			return true
		}

	default:
		// nop
	}

	idl := m.GetIdLike()

	if !Contains(idl, k) {
		return false
	}

	return true
}

func FormattedString(k KennungSansGattung) string {
	sb := &strings.Builder{}
	parts := k.Parts()
	sb.WriteString(parts[0])
	sb.WriteString(parts[1])
	sb.WriteString(parts[2])
	return sb.String()
}

func AlignedParts(
	id Kennung,
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

func Aligned(id Kennung, lenLeft, lenRight int) string {
	left, middle, right := AlignedParts(id, lenLeft, lenRight)
	return fmt.Sprintf("%s%s%s", left, middle, right)
}

func LeftSubtract[
	T schnittstellen.Stringer,
	TPtr schnittstellen.StringSetterPtr[T],
](
	a, b T,
) (c T, err error) {
	if err = TPtr(&c).Set(strings.TrimPrefix(a.String(), b.String())); err != nil {
		err = errors.Wrapf(err, "'%s' - '%s'", a, b)
		return
	}

	return
}

func ContainsWithoutUnderscoreSuffix[T schnittstellen.Stringer](a, b T) bool {
	as := []rune(a.String())
	bs := []rune(b.String())

	if len(bs) > len(as) {
		return false
	}

	if !strings.HasPrefix(a.String(), b.String()) {
		return false
	}

	if len(bs) == len(as) {
		return true
	}

	if as[len(bs)] == '_' {
		return false
	}

	return true
}

func ContainsExactly(a, b KennungSansGattung) bool {
	var (
		as = a.Parts()
		bs = b.Parts()
	)

	for i, e := range as {
		if bs[i] != e {
			return false
		}
	}

	return true
}

func Contains(a, b KennungSansGattung) bool {
	var (
		as = a.Parts()
		bs = b.Parts()
	)

	for i, e := range as {
		if !strings.HasPrefix(e, bs[i]) {
			return false
		}
	}

	return true
}

func Includes(a, b KennungSansGattung) bool {
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

	s3 = s4.CloneSetLike()

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

	s2 = s3.CloneSetLike()

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

	out = s1.CloneSetLike()

	return
}

func ExpandMany[T KennungLike[T], TPtr KennungLikePtr[T]](
	ks schnittstellen.Set[T],
	ex Expander,
) (out schnittstellen.Set[T]) {
	s1 := collections.MakeMutableSetStringer[T]()

	ks.Each(
		func(k T) (err error) {
			expandOne[T, TPtr](k, ex, s1)

			return
		},
	)

	out = s1.CloneSetLike()

	return
}

func Expanded(s EtikettSet, ex Expander) (out EtikettSet) {
	return ExpandMany[Etikett, *Etikett](s, ex)
}

func AddNormalized(es EtikettMutableSet, e Etikett) {
	ExpandOne(e, ExpanderRight).Each(es.Add)
	es.Add(e)

	c := es.CloneSetLike()
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
	s2 = s3.CloneSetLike()

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
		idx, hasMatch = BinarySearchForRuneInEtikettenSortedStringSlice(
			longer,
			c,
		)
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
