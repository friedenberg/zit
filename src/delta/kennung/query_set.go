package kennung

import (
	"encoding/gob"
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections"
)

const QueryNegationOperator rune = '^'

func init() {
	registerQuerySetGob[Etikett, *Etikett]()
}

func registerQuerySetGob[T QueryKennung[T], TPtr QueryKennungPtr[T]]() {
	gob.Register(querySet[T, TPtr]{})
	gob.Register(mutableQuerySet[T, TPtr]{})
	collections.RegisterGobTridexSet[T]()
}

type QueryKennung[T any] interface {
	KennungLike[T]
	QueryPrefixer
}

type QueryKennungPtr[T QueryKennung[T]] interface {
	KennungLikePtr[T]
}

type QuerySet[T QueryKennung[T], TPtr QueryKennungPtr[T]] interface {
	schnittstellen.Lenner
	schnittstellen.Stringer
	Contains(T) bool
	ContainsAgainst(schnittstellen.Set[T]) bool
	GetIncludes() schnittstellen.Set[T]
	GetExcludes() schnittstellen.Set[T]
	schnittstellen.ImmutableCloner[QuerySet[T, TPtr]]
	schnittstellen.MutableCloner[MutableQuerySet[T, TPtr]]
}

type MutableQuerySet[T QueryKennung[T], TPtr QueryKennungPtr[T]] interface {
	QuerySet[T, TPtr]
	AddString(string) error
	AddInclude(T) error
	AddExclude(T) error
}

func MakeMutableQuerySet[T QueryKennung[T], TPtr QueryKennungPtr[T]](
	ex func(string) (string, error),
	inc schnittstellen.Set[T],
	exc schnittstellen.Set[T],
) MutableQuerySet[T, TPtr] {
	if inc == nil {
		inc = collections.MakeSetStringer[T]()
	}

	if exc == nil {
		exc = collections.MakeSetStringer[T]()
	}

	return mutableQuerySet[T, TPtr]{
		Expander: ex,
		Include:  collections.MakeMutableTridexSet[T](inc.Elements()...),
		Exclude:  collections.MakeMutableTridexSet[T](exc.Elements()...),
	}
}

func MakeQuerySet[T QueryKennung[T], TPtr QueryKennungPtr[T]](
	ex func(string) (string, error),
	inc schnittstellen.Set[T],
	exc schnittstellen.Set[T],
) QuerySet[T, TPtr] {
	return querySet[T, TPtr]{
		set: MakeMutableQuerySet[T, TPtr](ex, inc, exc),
	}
}

type querySet[T QueryKennung[T], TPtr QueryKennungPtr[T]] struct {
	set MutableQuerySet[T, TPtr]
}

type mutableQuerySet[T QueryKennung[T], TPtr QueryKennungPtr[T]] struct {
	Expander func(string) (string, error)
	Include  collections.MutableTridexSet[T]
	Exclude  collections.MutableTridexSet[T]
}

//   ____  _        _
//  / ___|| |_ _ __(_)_ __   __ _  ___ _ __
//  \___ \| __| '__| | '_ \ / _` |/ _ \ '__|
//   ___) | |_| |  | | | | | (_| |  __/ |
//  |____/ \__|_|  |_|_| |_|\__, |\___|_|
//                          |___/

func (kqs querySet[T, TPtr]) String() string {
	return kqs.set.String()
}

func (kqs mutableQuerySet[T, TPtr]) String() string {
	sb := &strings.Builder{}

	var e T
	p := e.GetQueryPrefix()

	kqs.Include.Each(
		func(e T) (err error) {
			sb.WriteString(fmt.Sprintf("%s%s ", p, e))
			return
		},
	)

	kqs.Exclude.Each(
		func(e T) (err error) {
			sb.WriteString(fmt.Sprintf("%c%s%s ", QueryNegationOperator, p, e))
			return
		},
	)

	return sb.String()
}

//   _
//  | |    ___ _ __  _ __   ___ _ __
//  | |   / _ \ '_ \| '_ \ / _ \ '__|
//  | |__|  __/ | | | | | |  __/ |
//  |_____\___|_| |_|_| |_|\___|_|
//

func (kqs querySet[T, TPtr]) Len() int {
	return kqs.set.Len()
}

func (kqs mutableQuerySet[T, TPtr]) Len() int {
	return collections.Len(kqs.Include, kqs.Exclude)
}

//    ____            _        _
//   / ___|___  _ __ | |_ __ _(_)_ __  ___  ___ _ __
//  | |   / _ \| '_ \| __/ _` | | '_ \/ __|/ _ \ '__|
//  | |__| (_) | | | | || (_| | | | | \__ \  __/ |
//   \____\___/|_| |_|\__\__,_|_|_| |_|___/\___|_|
//

func (kqs querySet[T, TPtr]) ContainsAgainst(els schnittstellen.Set[T]) bool {
	return kqs.set.ContainsAgainst(els)
}

func (kqs mutableQuerySet[T, TPtr]) ContainsAgainst(els schnittstellen.Set[T]) bool {
	if els == nil || els.Len() == 0 {
		return kqs.Include.Len() == 0
	}

	if (kqs.Include.Len() == 0 || iter.Any(els, kqs.Include.Contains)) &&
		(kqs.Exclude.Len() == 0 || !iter.Any(els, kqs.Exclude.Contains)) {
		return true
	}

	return false
}

func (kqs querySet[T, TPtr]) Contains(e T) bool {
	return kqs.set.Contains(e)
}

func (kqs mutableQuerySet[T, TPtr]) Contains(e T) bool {
	if kqs.Include.Len() == 0 {
		return !kqs.Exclude.Contains(e)
	} else {
		return kqs.Include.Contains(e) && !kqs.Exclude.Contains(e)
	}
}

func (kqs querySet[T, TPtr]) GetIncludes() schnittstellen.Set[T] {
	return kqs.set.GetIncludes()
}

func (kqs querySet[T, TPtr]) GetExcludes() schnittstellen.Set[T] {
	return kqs.set.GetExcludes()
}

func (kqs mutableQuerySet[T, TPtr]) GetIncludes() schnittstellen.Set[T] {
	return kqs.Include.GetSet()
}

func (kqs mutableQuerySet[T, TPtr]) GetExcludes() schnittstellen.Set[T] {
	return kqs.Exclude.GetSet()
}

func (kqs querySet[T, TPtr]) ImmutableClone() QuerySet[T, TPtr] {
	return kqs.set.ImmutableClone()
}

func (kqs querySet[T, TPtr]) MutableClone() MutableQuerySet[T, TPtr] {
	return kqs.set.MutableClone()
}

func (kqs mutableQuerySet[T, TPtr]) ImmutableClone() QuerySet[T, TPtr] {
	return querySet[T, TPtr]{
		set: kqs.MutableClone(),
	}
}

func (kqs mutableQuerySet[T, TPtr]) MutableClone() MutableQuerySet[T, TPtr] {
	return mutableQuerySet[T, TPtr]{
		Expander: kqs.Expander,
		Include:  kqs.Include.MutableClone(),
		Exclude:  kqs.Exclude.MutableClone(),
	}
}

func (kqs mutableQuerySet[T, TPtr]) AddString(v string) (err error) {
	col := kqs.Include
	v = strings.TrimSpace(v)

	if len(v) > 0 && []rune(v)[0] == QueryNegationOperator {
		v = v[1:]
		col = kqs.Exclude
	}

	var e T
	p := e.GetQueryPrefix()

	if len(v) > 0 && v[:len(p)] == p {
		v = v[1:]
	}

	err = collections.ExpandAndAddString[T, TPtr](
		col,
		kqs.Expander,
		v,
	)

	return
}

func (kqs mutableQuerySet[T, TPtr]) AddInclude(e T) (err error) {
	return kqs.Include.Add(e)
}

func (kqs mutableQuerySet[T, TPtr]) AddExclude(e T) (err error) {
	return kqs.Exclude.Add(e)
}
