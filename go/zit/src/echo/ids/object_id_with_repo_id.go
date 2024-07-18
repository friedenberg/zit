package ids

import (
	"io"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
)

var poolObjectIdWithRepoId interfaces.Pool[ObjectIdWithRepoId, *ObjectIdWithRepoId]

func init() {
	poolObjectIdWithRepoId = pool.MakePool(
		nil,
		func(k *ObjectIdWithRepoId) {
			k.Reset()
		},
	)
}

func GetObjectIdWithRepoIdPool() interfaces.Pool[ObjectIdWithRepoId, *ObjectIdWithRepoId] {
	return poolObjectIdWithRepoId
}

type ObjectIdWithRepoId struct {
	ObjectId
	RepoId catgut.String
}

func MustObjectIdWithRepoId(kp IdLike) (k *ObjectIdWithRepoId) {
	k = &ObjectIdWithRepoId{}
	err := k.SetWithIdLike(kp)
	errors.PanicIfError(err)
	return
}

func (a *ObjectIdWithRepoId) GetRepoId() interfaces.RepoId {
	return MustRepoId(a.RepoId.String())
}

func (a *ObjectIdWithRepoId) Equals(b *ObjectIdWithRepoId) bool {
	if !a.ObjectId.Equals(&b.ObjectId) {
		return false
	}

	if !a.RepoId.Equals(&b.RepoId) {
		return false
	}

	return true
}

func (k3 *ObjectIdWithRepoId) WriteTo(w io.Writer) (n int64, err error) {
	err = todo.Implement()
	return
}

func (k3 *ObjectIdWithRepoId) ReadFrom(r io.Reader) (n int64, err error) {
	err = todo.Implement()
	return
}

func (k3 *ObjectIdWithRepoId) StringFromPtr() string {
	var sb strings.Builder

	sb.WriteRune('/')
	k3.RepoId.WriteTo(&sb)
	sb.WriteRune('/')
	sb.WriteString(k3.ObjectId.StringFromPtr())

	return sb.String()
}

func (k3 *ObjectIdWithRepoId) Len() int {
	return k3.RepoId.Len() + 1 + k3.ObjectId.Len()
}

func (k3 *ObjectIdWithRepoId) String() string {
	return k3.StringFromPtr()
}

func (k3 *ObjectIdWithRepoId) Reset() {
	k3.ObjectId.Reset()
	k3.RepoId.Reset()
}

func (k3 *ObjectIdWithRepoId) PartsStrings() IdParts {
	return IdParts{
		RepoId: &k3.RepoId,
		Left:   &k3.left,
		Middle: k3.middle,
		Right:  &k3.right,
	}
}

func (k3 *ObjectIdWithRepoId) GetGenre() interfaces.Genre {
	return k3.g
}

func MakeObjectIdWithRepoId(
	v interfaces.ObjectId,
	ka RepoId,
) (k *ObjectIdWithRepoId, err error) {
	k = &ObjectIdWithRepoId{}

	if err = k.RepoId.Set(ka.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = k.SetWithGenre(v.String(), v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (k3 *ObjectIdWithRepoId) Expand(
	a Abbr,
) (err error) {
	ex := a.ExpanderFor(k3.g)

	if ex == nil {
		return
	}

	v := k3.String()

	if v, err = ex(v); err != nil {
		err = nil
		return
	}

	if err = k3.SetWithGenre(v, k3.g); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (k3 *ObjectIdWithRepoId) Abbreviate(
	a Abbr,
) (err error) {
	return
}

func (h *ObjectIdWithRepoId) SetWithIdLike(
	k IdLike,
) (err error) {
	switch kt := k.(type) {
	case *ObjectIdWithRepoId:
		if err = kt.left.CopyTo(&h.left); err != nil {
			err = errors.Wrap(err)
			return
		}

		h.middle = kt.middle

		if err = kt.right.CopyTo(&h.right); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		p := k.Parts()

		if err = h.left.Set(p[0]); err != nil {
			err = errors.Wrap(err)
			return
		}

		mid := []byte(p[1])

		if len(mid) >= 1 {
			h.middle = mid[0]
		}

		if err = h.right.Set(p[2]); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	h.SetGenre(k)

	return
}

func (h *ObjectIdWithRepoId) SetWithGenre(
	v string,
	g interfaces.GenreGetter,
) (err error) {
	h.g = genres.Make(g.GetGenre())

	if err = h.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (h *ObjectIdWithRepoId) TodoSetBytes(v *catgut.String) (err error) {
	return h.Set(v.String())
}

func (h *ObjectIdWithRepoId) Set(v string) (err error) {
	var k IdLike

	switch h.g {
	case genres.Unknown:
		k, err = Make(v)

	case genres.Zettel:
		var h ZettelId
		err = h.Set(v)
		k = h

	case genres.Tag:
		var h Tag
		err = h.Set(v)
		k = h

	case genres.Type:
		var h Type
		err = h.Set(v)
		k = h

	case genres.Repo:
		var h RepoId
		err = h.Set(v)
		k = h

	case genres.Config:
		var h Config
		err = h.Set(v)
		k = h

	case genres.InventoryList:
		var h Tai
		err = h.Set(v)
		k = h

	default:
		err = genres.MakeErrUnrecognizedGenre(h.g.GetGenreString())
	}

	if err != nil {
		err = errors.Wrapf(err, "String: %q", v)
		return
	}

	if err = h.SetWithIdLike(k); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *ObjectIdWithRepoId) ResetWith(b *ObjectIdWithRepoId) {
	a.g = b.g
	b.left.CopyTo(&a.left)
	b.right.CopyTo(&a.right)
	a.middle = b.middle
}

func (a *ObjectIdWithRepoId) ResetWithObjectId(b IdLike) (err error) {
	return a.SetWithIdLike(b)
}

func (t *ObjectIdWithRepoId) MarshalText() (text []byte, err error) {
	text = []byte(FormattedString(t))
	return
}

func (t *ObjectIdWithRepoId) UnmarshalText(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (t *ObjectIdWithRepoId) MarshalBinary() (text []byte, err error) {
	text = []byte(FormattedString(t))
	return
}

func (t *ObjectIdWithRepoId) UnmarshalBinary(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
