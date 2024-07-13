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

type IdWithoutGenre interface {
	interfaces.Stringer
	Parts() [3]string
}

type Id interface {
	IdWithoutGenre
	interfaces.GenreGetter
}

type KennungSansGattungPtr interface {
	IdWithoutGenre
	interfaces.Resetter
	interfaces.Setter
}

type KennungPtr interface {
	Id
	KennungSansGattungPtr
}

type KennungLike[T any] interface {
	Id
	interfaces.GenreGetter
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
		var h Config

		if err = h.Set(v); err == nil {
			k = &h
			return
		}
	}

	{
		var e Tag

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
		var ka RepoId

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

func Equals(a, b Id) (ok bool) {
	if a.GetGenre().GetGenreString() != b.GetGenre().GetGenreString() {
		return
	}

	if a.String() != b.String() {
		return
	}

	return true
}

func FormattedString(k IdWithoutGenre) string {
	sb := &strings.Builder{}
	parts := k.Parts()
	sb.WriteString(parts[0])
	sb.WriteString(parts[1])
	sb.WriteString(parts[2])
	return sb.String()
}

func AlignedParts(
	id Id,
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

func Aligned(id Id, lenLeft, lenRight int) string {
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

func ContainsExactly(a, b IdWithoutGenre) bool {
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

func Contains(a, b IdWithoutGenre) bool {
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

func Includes(a, b IdWithoutGenre) bool {
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

func SansPrefix(a Tag) (b Tag) {
	b = MustTag(strings.TrimPrefix(a.String(), "-"))
	return
}

func IsDependentLeaf(a Tag) (has bool) {
	has = strings.HasPrefix(strings.TrimSpace(a.String()), "-")
	return
}

func HasParentPrefix(a, b Tag) (has bool) {
	has = strings.HasPrefix(strings.TrimSpace(a.String()), b.String())
	return
}

func IntersectPrefixes(haystack TagSet, needle Tag) (s3 TagSet) {
	s4 := MakeTagMutableSet()

	for _, e := range iter.Elements[Tag](haystack) {
		if strings.HasPrefix(e.String(), needle.String()) {
			s4.Add(e)
		}
	}

	s3 = s4.CloneSetPtrLike()

	return
}

func SubtractPrefix(s1 TagSet, e Tag) (s2 TagSet) {
	s3 := MakeTagMutableSet()

	for _, e1 := range iter.Elements[Tag](s1) {
		e2, _ := LeftSubtract(e1, e)

		if e2.String() == "" {
			continue
		}

		s3.Add(e2)
	}

	s2 = s3.CloneSetPtrLike()

	return
}

func Description(s TagSet) string {
	return iter.StringCommaSeparated[Tag](s)
}

func WithRemovedCommonPrefixes(s TagSet) (s2 TagSet) {
	es1 := iter.SortedValues[Tag](s)
	es := make([]Tag, 0, len(es1))

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

	s2 = MakeTagSet(es...)

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

func Expanded(s TagSet, ex expansion.Expander) (out TagSet) {
	return ExpandMany(s, ex)
}

func AddNormalizedEtikett(es TagMutableSet, e *Tag) {
	ExpandOne(e, expansion.ExpanderRight).Each(es.Add)
	errors.PanicIfError(iter.AddClonePool(
		es,
		GetTagPool(),
		TagResetter,
		e,
	))

	c := es.CloneSetPtrLike()
	es.Reset()
	WithRemovedCommonPrefixes(c).Each(es.Add)
}

func RemovePrefixes(es TagMutableSet, needle Tag) {
	for _, haystack := range iter.Elements(es) {
		// TODO-P2 make more efficient
		if strings.HasPrefix(haystack.String(), needle.String()) {
			es.Del(haystack)
		}
	}
}

func Withdraw(s1 TagMutableSet, e Tag) (s2 TagSet) {
	s3 := MakeTagMutableSet()

	for _, e1 := range iter.Elements[Tag](s1) {
		if Contains(e1, e) {
			s3.Add(e1)
		}
	}

	s3.Each(s1.Del)
	s2 = s3.CloneSetPtrLike()

	return
}
