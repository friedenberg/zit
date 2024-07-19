package ids

import (
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
)

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

func AddNormalizedTag(es TagMutableSet, e *Tag) {
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
