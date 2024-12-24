package ids

import (
	"fmt"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type Abbreviatable interface {
	interfaces.Stringer
}

type IdWithParts interface {
	interfaces.Stringer
	Parts() [3]string
}

type (
	IdLike interface {
		IdWithParts
		interfaces.GenreGetter
		// interfaces.Resetter
		// interfaces.Setter
	}

	ObjectIdLike interface {
		interfaces.GenreGetter
		interfaces.Stringer
		// GetObjectId() *ObjectId
		IsEmpty() bool
	}

	ExternalObjectId interface {
		ObjectIdLike
		interfaces.GenreGetter
		interfaces.Stringer
		ExternalObjectIdGetter
		ExternalObjectIdCloner
	}

	ExternalObjectIdGetter interface {
		GetExternalObjectId() ExternalObjectId
	}

	ExternalObjectIdCloner interface {
		CloneExternalObjectId() ExternalObjectId
	}
)

type Index struct{}

func Make(v string) (k IdLike, err error) {
	if v == "" {
		k = &ObjectId{}
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
		var h Tai

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
		var t Type

		if err = t.Set(v); err == nil {
			k = &t
			return
		}
	}

	{
		var h ZettelId

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

	err = errors.Errorf("%q is not a valid object id", v)

	return
}

func Equals(a, b IdLike) (ok bool) {
	if a.GetGenre().GetGenreString() != b.GetGenre().GetGenreString() {
		return
	}

	if a.String() != b.String() {
		return
	}

	return true
}

func FormattedString(k IdWithParts) string {
	sb := &strings.Builder{}
	parts := k.Parts()
	sb.WriteString(parts[0])
	sb.WriteString(parts[1])
	sb.WriteString(parts[2])
	return sb.String()
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

func Contains(a, b IdWithParts) bool {
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

func ContainsExactly(a, b IdWithParts) bool {
	var (
		as = a.Parts()
		bs = b.Parts()
	)

	if as != bs {
		return false
	}

	return true
}

func IsEmpty[T interfaces.Stringer](a T) bool {
	return len(a.String()) == 0
}
