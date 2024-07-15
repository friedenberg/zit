package ids

import (
	"io"
	"math"
	"slices"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_extensions"
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
	g                   genres.Genre
	middle              byte // remove and replace with virtual
	kasten, left, right catgut.String
}

func MustObjectIdWithRepoId(kp IdLike) (k *ObjectIdWithRepoId) {
	k = &ObjectIdWithRepoId{}
	err := k.SetWithIdLike(kp)
	errors.PanicIfError(err)
	return
}

func (a *ObjectIdWithRepoId) GetRepoId() interfaces.RepoId {
	return MustRepoId(a.kasten.String())
}

func (a *ObjectIdWithRepoId) IsVirtual() bool {
	switch a.g {
	case genres.Zettel:
		return slices.Equal(a.left.Bytes(), []byte{'%'})

	case genres.Tag:
		return a.middle == '%' || slices.Equal(a.left.Bytes(), []byte{'%'})

	default:
		return false
	}
}

func (a *ObjectIdWithRepoId) Equals(b *ObjectIdWithRepoId) bool {
	if a.g != b.g {
		return false
	}

	if a.middle != b.middle {
		return false
	}

	if !a.left.Equals(&b.left) {
		return false
	}

	if !a.right.Equals(&b.right) {
		return false
	}

	if !a.kasten.Equals(&b.kasten) {
		return false
	}

	return true
}

func (k3 *ObjectIdWithRepoId) WriteTo(w io.Writer) (n int64, err error) {
	if k3.Len() > math.MaxUint8 {
		err = errors.Errorf(
			"%q is greater than max uint8 (%d)",
			k3.String(),
			math.MaxUint8,
		)

		return
	}

	var n1 int64
	n1, err = k3.g.WriteTo(w)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	b := [2]uint8{uint8(k3.Len()), uint8(k3.left.Len())}

	var n2 int
	n2, err = ohio.WriteAllOrDieTrying(w, b[:])
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = k3.left.WriteTo(w)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	bMid := [1]byte{k3.middle}

	n2, err = ohio.WriteAllOrDieTrying(w, bMid[:])
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = k3.right.WriteTo(w)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (k3 *ObjectIdWithRepoId) ReadFrom(r io.Reader) (n int64, err error) {
	var n1 int64
	n1, err = k3.g.ReadFrom(r)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	var b [2]uint8

	var n2 int
	n2, err = ohio.ReadAllOrDieTrying(r, b[:])
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	contentLength := b[0]
	middlePos := b[1]

	if middlePos > contentLength-1 {
		err = errors.Errorf(
			"middle position %d is greater than last index: %d",
			middlePos,
			contentLength,
		)
		return
	}

	if _, err = k3.left.ReadNFrom(r, int(middlePos)); err != nil {
		err = errors.Wrap(err)
		return
	}

	var bMiddle [1]uint8

	n2, err = ohio.ReadAllOrDieTrying(r, bMiddle[:])
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	k3.middle = bMiddle[0]

	if _, err = k3.right.ReadNFrom(r, int(contentLength-middlePos-1)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (k3 *ObjectIdWithRepoId) SetGattung(g interfaces.GenreGetter) {
	if g == nil {
		k3.g = genres.Unknown
	} else {
		k3.g = genres.Must(g.GetGenre())
	}

	if k3.g == genres.Zettel {
		k3.middle = '/'
	}
}

func (k3 *ObjectIdWithRepoId) StringFromPtr() string {
	var sb strings.Builder

	switch k3.g {
	case genres.Zettel:
		sb.Write(k3.left.Bytes())
		sb.WriteByte(k3.middle)
		sb.Write(k3.right.Bytes())

	case genres.Type:
		sb.Write(k3.right.Bytes())

	default:
		if k3.left.Len() > 0 {
			sb.Write(k3.left.Bytes())
		}

		if k3.middle != '\x00' {
			sb.WriteByte(k3.middle)
		}

		if k3.right.Len() > 0 {
			sb.Write(k3.right.Bytes())
		}
	}

	return sb.String()
}

func (k3 *ObjectIdWithRepoId) IsEmpty() bool {
	if k3.g == genres.Zettel {
		if k3.left.IsEmpty() && k3.right.IsEmpty() {
			return true
		}
	}

	return k3.left.Len() == 0 && k3.middle == 0 && k3.right.Len() == 0
}

func (k3 *ObjectIdWithRepoId) Len() int {
	return k3.left.Len() + 1 + k3.right.Len()
}

func (k3 *ObjectIdWithRepoId) KopfUndSchwanz() (kopf, schwanz string) {
	kopf = k3.left.String()
	schwanz = k3.right.String()

	return
}

func (k3 *ObjectIdWithRepoId) LenKopfUndSchwanz() (int, int) {
	return k3.left.Len(), k3.right.Len()
}

func (k3 *ObjectIdWithRepoId) String() string {
	return k3.StringFromPtr()
}

func (k3 *ObjectIdWithRepoId) Reset() {
	k3.g = genres.Unknown
	k3.left.Reset()
	k3.middle = 0
	k3.right.Reset()
}

func (k3 *ObjectIdWithRepoId) PartsStrings() IdParts {
	return IdParts{
		RepoId: &k3.kasten,
		Left:   &k3.left,
		Middle: k3.middle,
		Right:  &k3.right,
	}
}

func (k3 *ObjectIdWithRepoId) Parts() [3]string {
	var mid string

	if k3.middle != 0 {
		mid = string([]byte{k3.middle})
	}

	return [3]string{
		k3.left.String(),
		mid,
		k3.right.String(),
	}
}

func (k3 *ObjectIdWithRepoId) GetGenre() interfaces.Genre {
	return k3.g
}

func MakeObjectIdWithRepoId(
	v interfaces.ObjectId,
	ka RepoId,
) (k *ObjectIdWithRepoId, err error) {
	k = &ObjectIdWithRepoId{
		g: genres.Unknown,
	}

	if err = k.kasten.Set(ka.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = k.SetWithGattung(v.String(), v); err != nil {
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

	if err = k3.SetWithGattung(v, k3.g); err != nil {
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

func (k3 *ObjectIdWithRepoId) SetFromPath(
	path string,
	fe file_extensions.FileExtensions,
) (err error) {
	els := files.PathElements(path)
	ext := els[0]

	switch ext {
	case fe.Etikett:
		if err = k3.SetWithGattung(els[1], genres.Tag); err != nil {
			err = errors.Wrap(err)
			return
		}

	case fe.Typ:
		if err = k3.SetWithGattung(els[1], genres.Type); err != nil {
			err = errors.Wrap(err)
			return
		}

	case fe.Kasten:
		if err = k3.SetWithGattung(els[1], genres.Repo); err != nil {
			err = errors.Wrap(err)
			return
		}

	case fe.Zettel:
		if err = k3.SetWithGattung(els[2]+"/"+els[1], genres.Zettel); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		err = ErrFDNotId
		return
	}

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

	h.SetGattung(k)

	return
}

func (h *ObjectIdWithRepoId) SetWithGattung(
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

func (a *ObjectIdWithRepoId) ResetWithKennung(b IdLike) (err error) {
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
	// if t.g == gattung.Unknown {
	// 	err = errors.Wrapf(gattung.ErrEmptyKennung{}, "Kennung: %s", t)
	// 	return
	// }

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
