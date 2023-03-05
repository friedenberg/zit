package kennung

import (
	"encoding/gob"
	"strings"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
)

func init() {
	registerQuerySetGob[Etikett, *Etikett]()
}

func registerQuerySetGob[T KennungLike[T], TPtr KennungLikePtr[T]]() {
	gob.Register(querySet[T, TPtr]{})
	gob.Register(mutableQuerySet[T, TPtr]{})
}

type QuerySet[T KennungLike[T], TPtr KennungLikePtr[T]] interface {
	schnittstellen.Lenner
	schnittstellen.Stringer
	Contains(T) bool
	GetIncludes() schnittstellen.Set[T]
	GetExcludes() schnittstellen.Set[T]
	schnittstellen.ImmutableCloner[QuerySet[T, TPtr]]
	schnittstellen.MutableCloner[MutableQuerySet[T, TPtr]]
}

type MutableQuerySet[T KennungLike[T], TPtr KennungLikePtr[T]] interface {
	QuerySet[T, TPtr]
	AddString(string) error
	AddInclude(T) error
	AddExclude(T) error
	// schnittstellen.Adder[T]
}

func MakeMutableQuerySet[T KennungLike[T], TPtr KennungLikePtr[T]](
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
		Include:  inc.MutableClone(),
		Exclude:  exc.MutableClone(),
	}
}

func MakeQuerySet[T KennungLike[T], TPtr KennungLikePtr[T]](
	ex func(string) (string, error),
	inc schnittstellen.Set[T],
	exc schnittstellen.Set[T],
) QuerySet[T, TPtr] {
	if inc == nil {
		inc = collections.MakeSetStringer[T]()
	}

	if exc == nil {
		exc = collections.MakeSetStringer[T]()
	}

	return querySet[T, TPtr]{
		Expander: ex,
		Include:  inc,
		Exclude:  exc,
	}
}

type querySet[T KennungLike[T], TPtr KennungLikePtr[T]] struct {
	Expander func(string) (string, error)
	Include  schnittstellen.Set[T]
	Exclude  schnittstellen.Set[T]
}

type mutableQuerySet[T KennungLike[T], TPtr KennungLikePtr[T]] struct {
	Expander func(string) (string, error)
	Include  schnittstellen.MutableSet[T]
	Exclude  schnittstellen.MutableSet[T]
}

//   ____  _        _
//  / ___|| |_ _ __(_)_ __   __ _  ___ _ __
//  \___ \| __| '__| | '_ \ / _` |/ _ \ '__|
//   ___) | |_| |  | | | | | (_| |  __/ |
//  |____/ \__|_|  |_|_| |_|\__, |\___|_|
//                          |___/

func (kqs querySet[T, TPtr]) String() string {
	return ""
}

func (kqs mutableQuerySet[T, TPtr]) String() string {
	return ""
}

//   _
//  | |    ___ _ __  _ __   ___ _ __
//  | |   / _ \ '_ \| '_ \ / _ \ '__|
//  | |__|  __/ | | | | | |  __/ |
//  |_____\___|_| |_|_| |_|\___|_|
//

func (kqs querySet[T, TPtr]) Len() int {
	return collections.Len(kqs.Include, kqs.Exclude)
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

func (kqs querySet[T, TPtr]) Contains(e T) bool {
	return kqs.Include.Contains(e) && !kqs.Exclude.Contains(e)
}

func (kqs mutableQuerySet[T, TPtr]) Contains(e T) bool {
	return kqs.Include.Contains(e) && !kqs.Exclude.Contains(e)
}

func (kqs querySet[T, TPtr]) GetIncludes() schnittstellen.Set[T] {
	return kqs.Include
}

func (kqs querySet[T, TPtr]) GetExcludes() schnittstellen.Set[T] {
	return kqs.Exclude
}

func (kqs mutableQuerySet[T, TPtr]) GetIncludes() schnittstellen.Set[T] {
	return kqs.Include
}

func (kqs mutableQuerySet[T, TPtr]) GetExcludes() schnittstellen.Set[T] {
	return kqs.Exclude
}

func (kqs querySet[T, TPtr]) ImmutableClone() QuerySet[T, TPtr] {
	return querySet[T, TPtr]{
		Expander: kqs.Expander,
		Include:  kqs.Include,
		Exclude:  kqs.Exclude,
	}
}

func (kqs querySet[T, TPtr]) MutableClone() MutableQuerySet[T, TPtr] {
	return mutableQuerySet[T, TPtr]{
		Expander: kqs.Expander,
		Include:  kqs.Include.MutableClone(),
		Exclude:  kqs.Exclude.MutableClone(),
	}
}

func (kqs mutableQuerySet[T, TPtr]) ImmutableClone() QuerySet[T, TPtr] {
	return querySet[T, TPtr]{
		Expander: kqs.Expander,
		Include:  kqs.Include,
		Exclude:  kqs.Exclude,
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

	if len(v) > 0 && v[0] == '!' {
		v = v[1:]
		col = kqs.Exclude
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
