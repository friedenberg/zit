package ids

import (
	"sort"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_ptr"
)

type (
	TagSet        = interfaces.SetPtrLike[Tag, *Tag]
	TagMutableSet = interfaces.MutableSetPtrLike[Tag, *Tag]
)

var TagSetEmpty TagSet

func init() {
	collections_ptr.RegisterGobValue[Tag](nil)
	TagSetEmpty = MakeTagSet()
}

func MakeTagSet(es ...Tag) (s TagSet) {
	if len(es) == 0 && TagSetEmpty != nil {
		return TagSetEmpty
	}

	return TagSet(
		collections_ptr.MakeValueSetValue(nil, es...),
	)
}

func MakeTagSetStrings(vs ...string) (s TagSet, err error) {
	return collections_ptr.MakeValueSetString[Tag](nil, vs...)
}

func MakeMutableTagSet(hs ...Tag) TagMutableSet {
	return MakeTagMutableSet(hs...)
}

func MakeTagMutableSet(hs ...Tag) TagMutableSet {
	return TagMutableSet(
		collections_ptr.MakeMutableValueSetValue(
			nil,
			hs...,
		),
	)
}

func TagSetEquals(a, b TagSet) bool {
	return quiter.SetEqualsPtr(a, b)
}

type TagSlice []Tag

func MakeTagSlice(es ...Tag) (s TagSlice) {
	s = make([]Tag, len(es))

	for i, e := range es {
		s[i] = e
	}

	return
}

func NewSliceFromStrings(es ...string) (s TagSlice, err error) {
	s = make([]Tag, len(es))

	for i, e := range es {
		if err = s[i].Set(e); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *TagSlice) DropFirst() {
	if s.Len() > 0 {
		*s = (*s)[1:]
	}
}

func (s TagSlice) Len() int {
	return len(s)
}

func (es *TagSlice) AddString(v string) (err error) {
	var e Tag

	if err = e.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	es.Add(e)

	return
}

func (es *TagSlice) Add(e Tag) {
	*es = append(*es, e)
}

func (s *TagSlice) Set(v string) (err error) {
	es := strings.Split(v, ",")

	for _, e := range es {
		if err = s.AddString(e); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (es TagSlice) SortedString() (out []string) {
	out = make([]string, len(es))

	i := 0

	for _, e := range es {
		out[i] = e.String()
		i++
	}

	sort.Slice(
		out,
		func(i, j int) bool {
			return out[i] < out[j]
		},
	)

	return
}

func (s TagSlice) String() string {
	return strings.Join(s.SortedString(), ", ")
}

func (s TagSlice) ToSet() TagSet {
	return MakeTagSet([]Tag(s)...)
}
