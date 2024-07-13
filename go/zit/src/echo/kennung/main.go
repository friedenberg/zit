package kennung

import (
	"fmt"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_ptr"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
)

type QueryPrefixer interface {
	GetQueryPrefix() string
}

type KennungSansGattung interface {
	interfaces.Stringer
	Parts() [3]string
}

type Kennung interface {
	KennungSansGattung
	interfaces.GattungGetter
}

type KennungSansGattungPtr interface {
	KennungSansGattung
	interfaces.Resetter
	interfaces.Setter
}

type KennungPtr interface {
	Kennung
	KennungSansGattungPtr
}

type KennungLike[T any] interface {
	Kennung
	interfaces.GattungGetter
	interfaces.Stringer
}

type KennungLikePtr[T KennungLike[T]] interface {
	interfaces.Ptr[T]
	KennungLike[T]
	KennungPtr
	interfaces.SetterPtr[T]
}

type Index struct{}

func Make(v string) (k KennungPtr, err error) {
	if v == "" {
		k = &Kennung2{}
		return
	}

	{
		var h Konfig

		if err = h.Set(v); err == nil {
			k = &h
			return
		}
	}

	{
		var e Etikett

		if err = e.Set(v); err == nil {
			k = &e
			return
		}
	}

	{
		var t Typ

		if err = t.Set(v); err == nil {
			k = &t
			return
		}
	}

	{
		var h Hinweis

		if err = h.Set(v); err == nil {
			k = &h
			return
		}
	}

	{
		var ka Kasten

		if err = ka.Set(v); err == nil {
			k = &ka
			return
		}
	}

	{
		var h Tai

		if err = h.Set(v); err == nil {
			k = &h
			return
		}
	}

	err = errors.Errorf("%q is not a valid Kennung", v)

	return
}

func Equals(a, b Kennung) (ok bool) {
	if a.GetGattung().GetGattungString() != b.GetGattung().GetGattungString() {
		return
	}

	if a.String() != b.String() {
		return
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
	T interfaces.Stringer,
	TPtr interfaces.StringSetterPtr[T],
](
	a, b T,
) (c T, err error) {
	if err = TPtr(&c).Set(strings.TrimPrefix(a.String(), b.String())); err != nil {
		err = errors.Wrapf(err, "'%s' - '%s'", a, b)
		return
	}

	return
}

func ContainsWithoutUnderscoreSuffix[T interfaces.Stringer](a, b T) bool {
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

func Less[T interfaces.Stringer](a, b T) bool {
	return a.String() < b.String()
}

func LessLen[T interfaces.Stringer](a, b T) bool {
	return len(a.String()) < len(b.String())
}

func IsEmpty[T interfaces.Stringer](a T) bool {
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

func IntersectPrefixes(haystack EtikettSet, needle Etikett) (s3 EtikettSet) {
	s4 := MakeEtikettMutableSet()

	for _, e := range iter.Elements[Etikett](haystack) {
		if strings.HasPrefix(e.String(), needle.String()) {
			s4.Add(e)
		}
	}

	s3 = s4.CloneSetPtrLike()

	return
}

func SubtractPrefix(s1 EtikettSet, e Etikett) (s2 EtikettSet) {
	s3 := MakeEtikettMutableSet()

	for _, e1 := range iter.Elements[Etikett](s1) {
		e2, _ := LeftSubtract(e1, e)

		if e2.String() == "" {
			continue
		}

		s3.Add(e2)
	}

	s2 = s3.CloneSetPtrLike()

	return
}

func Description(s EtikettSet) string {
	return iter.StringCommaSeparated[Etikett](s)
}

func WithRemovedCommonPrefixes(s EtikettSet) (s2 EtikettSet) {
	es1 := iter.SortedValues[Etikett](s)
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
	k TPtr,
	ex expansion.Expander,
	acc interfaces.Adder[T],
) {
	f := iter.MakeFuncSetString[T, TPtr](acc)
	ex.Expand(f, k.String())
}

func ExpandOneSlice[T KennungLike[T], TPtr KennungLikePtr[T]](
	k TPtr,
	exes ...expansion.Expander,
) (out []T) {
	s1 := collections_value.MakeMutableValueSet[T](nil)

	if len(exes) == 0 {
		exes = []expansion.Expander{expansion.ExpanderAll}
	}

	for _, ex := range exes {
		expandOne(k, ex, s1)
	}

	out = iter.SortedValuesBy(
		s1,
		func(a, b T) bool {
			return len(a.String()) < len(b.String())
		},
	)

	return
}

func ExpandOne[T KennungLike[T], TPtr KennungLikePtr[T]](
	k TPtr,
	exes ...expansion.Expander,
) (out interfaces.SetPtrLike[T, TPtr]) {
	s1 := collections_ptr.MakeMutableValueSetValue[T, TPtr](nil)

	if len(exes) == 0 {
		exes = []expansion.Expander{expansion.ExpanderAll}
	}

	for _, ex := range exes {
		expandOne(k, ex, s1)
	}

	out = s1.CloneSetPtrLike()

	return
}

func ExpandMany[T KennungLike[T], TPtr KennungLikePtr[T]](
	ks interfaces.SetPtrLike[T, TPtr],
	ex expansion.Expander,
) (out interfaces.SetPtrLike[T, TPtr]) {
	s1 := collections_ptr.MakeMutableValueSetValue[T, TPtr](nil)

	ks.EachPtr(
		func(k TPtr) (err error) {
			expandOne[T, TPtr](k, ex, s1)

			return
		},
	)

	out = s1.CloneSetPtrLike()

	return
}

func ExpandOneTo[T KennungLike[T], TPtr KennungLikePtr[T]](
	k TPtr,
	ex expansion.Expander,
	s1 interfaces.FuncSetString,
) (out interfaces.SetPtrLike[T, TPtr]) {
	ex.Expand(s1, k.String())

	return
}

func Expanded(s EtikettSet, ex expansion.Expander) (out EtikettSet) {
	return ExpandMany(s, ex)
}

func AddNormalizedEtikett(es EtikettMutableSet, e *Etikett) {
	ExpandOne(e, expansion.ExpanderRight).Each(es.Add)
	errors.PanicIfError(iter.AddClonePool(
		es,
		GetEtikettPool(),
		EtikettResetter,
		e,
	))

	c := es.CloneSetPtrLike()
	es.Reset()
	WithRemovedCommonPrefixes(c).Each(es.Add)
}

func RemovePrefixes(es EtikettMutableSet, needle Etikett) {
	for _, haystack := range iter.Elements(es) {
		// TODO-P2 make more efficient
		if strings.HasPrefix(haystack.String(), needle.String()) {
			es.Del(haystack)
		}
	}
}

func Withdraw(s1 EtikettMutableSet, e Etikett) (s2 EtikettSet) {
	s3 := MakeEtikettMutableSet()

	for _, e1 := range iter.Elements[Etikett](s1) {
		if Contains(e1, e) {
			s3.Add(e1)
		}
	}

	s3.Each(s1.Del)
	s2 = s3.CloneSetPtrLike()

	return
}
