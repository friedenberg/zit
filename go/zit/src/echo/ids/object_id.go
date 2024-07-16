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

var poolObjectId interfaces.Pool[ObjectId, *ObjectId]

func init() {
	poolObjectId = pool.MakePool(
		nil,
		func(k *ObjectId) {
			k.Reset()
		},
	)
}

func GetObjectIdPool() interfaces.Pool[ObjectId, *ObjectId] {
	return poolObjectId
}

type ObjectId struct {
	g           genres.Genre
	middle      byte // remove and replace with virtual
	left, right catgut.String
}

func (a *ObjectId) Clone() (b *ObjectId) {
	b = GetObjectIdPool().Get()
	b.ResetWithIdLike(a)
	return
}

func MustObjectId(kp IdLike) (k *ObjectId) {
	k = &ObjectId{}
	err := k.SetWithIdLike(kp)
	errors.PanicIfError(err)
	return
}

func (a *ObjectId) IsVirtual() bool {
	switch a.g {
	case genres.Zettel:
		return slices.Equal(a.left.Bytes(), []byte{'%'})

	case genres.Tag:
		return a.middle == '%' || slices.Equal(a.left.Bytes(), []byte{'%'})

	default:
		return false
	}
}

func (a *ObjectId) Equals(b *ObjectId) bool {
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

	return true
}

func (k2 *ObjectId) WriteTo(w io.Writer) (n int64, err error) {
	if k2.Len() > math.MaxUint8 {
		err = errors.Errorf(
			"%q is greater than max uint8 (%d)",
			k2.String(),
			math.MaxUint8,
		)

		return
	}

	var n1 int64
	n1, err = k2.g.WriteTo(w)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	b := [2]uint8{uint8(k2.Len()), uint8(k2.left.Len())}

	var n2 int
	n2, err = ohio.WriteAllOrDieTrying(w, b[:])
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = k2.left.WriteTo(w)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	bMid := [1]byte{k2.middle}

	n2, err = ohio.WriteAllOrDieTrying(w, bMid[:])
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = k2.right.WriteTo(w)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (k2 *ObjectId) ReadFrom(r io.Reader) (n int64, err error) {
	var n1 int64
	n1, err = k2.g.ReadFrom(r)
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

	if _, err = k2.left.ReadNFrom(r, int(middlePos)); err != nil {
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

	k2.middle = bMiddle[0]

	if _, err = k2.right.ReadNFrom(r, int(contentLength-middlePos-1)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (k2 *ObjectId) SetGenre(g interfaces.GenreGetter) {
	if g == nil {
		k2.g = genres.Unknown
	} else {
		k2.g = genres.Must(g.GetGenre())
	}

	if k2.g == genres.Zettel {
		k2.middle = '/'
	}
}

func (k2 *ObjectId) StringFromPtr() string {
	var sb strings.Builder

	switch k2.g {
	case genres.Zettel:
		sb.Write(k2.left.Bytes())
		sb.WriteByte(k2.middle)
		sb.Write(k2.right.Bytes())

	case genres.Type:
		sb.Write(k2.right.Bytes())

	default:
		if k2.left.Len() > 0 {
			sb.Write(k2.left.Bytes())
		}

		if k2.middle != '\x00' {
			sb.WriteByte(k2.middle)
		}

		if k2.right.Len() > 0 {
			sb.Write(k2.right.Bytes())
		}
	}

	return sb.String()
}

func (k2 *ObjectId) IsEmpty() bool {
	if k2.g == genres.Zettel {
		if k2.left.IsEmpty() && k2.right.IsEmpty() {
			return true
		}
	}

	return k2.left.Len() == 0 && k2.middle == 0 && k2.right.Len() == 0
}

func (k2 *ObjectId) Len() int {
	return k2.left.Len() + 1 + k2.right.Len()
}

func (k2 *ObjectId) GetHeadAndTail() (head, tail string) {
	head = k2.left.String()
	tail = k2.right.String()

	return
}

func (k2 *ObjectId) LenHeadAndTail() (int, int) {
	return k2.left.Len(), k2.right.Len()
}

func (k2 *ObjectId) String() string {
	return k2.StringFromPtr()
}

func (k2 *ObjectId) Reset() {
	k2.g = genres.Unknown
	k2.left.Reset()
	k2.middle = 0
	k2.right.Reset()
}

type IdParts struct {
	Middle              byte
	RepoId, Left, Right *catgut.String
}

func (k2 *ObjectId) PartsStrings() IdParts {
	return IdParts{
		Left:   &k2.left,
		Middle: k2.middle,
		Right:  &k2.right,
	}
}

func (k2 *ObjectId) Parts() [3]string {
	var mid string

	if k2.middle != 0 {
		mid = string([]byte{k2.middle})
	}

	return [3]string{
		k2.left.String(),
		mid,
		k2.right.String(),
	}
}

func (k2 *ObjectId) GetGenre() interfaces.Genre {
	return k2.g
}

func MakeId(v string) (IdLikePtr, error) {
	k := &ObjectId{
		g: genres.Unknown,
	}

	return k, k.Set(v)
}

func (k2 *ObjectId) Expand(
	a Abbr,
) (err error) {
	ex := a.ExpanderFor(k2.g)

	if ex == nil {
		return
	}

	v := k2.String()

	if v, err = ex(v); err != nil {
		err = nil
		return
	}

	if err = k2.SetWithGenre(v, k2.g); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (k2 *ObjectId) Abbreviate(
	a Abbr,
) (err error) {
	return
}

var ErrFDNotId = errors.New("not a id file")

func (k2 *ObjectId) SetFromPath(
	path string,
	fe file_extensions.FileExtensions,
) (err error) {
	els := files.PathElements(path)
	ext := els[0]

	switch ext {
	case fe.Etikett:
		if err = k2.SetWithGenre(els[1], genres.Tag); err != nil {
			err = errors.Wrap(err)
			return
		}

	case fe.Typ:
		if err = k2.SetWithGenre(els[1], genres.Type); err != nil {
			err = errors.Wrap(err)
			return
		}

	case fe.Kasten:
		if err = k2.SetWithGenre(els[1], genres.Repo); err != nil {
			err = errors.Wrap(err)
			return
		}

	case fe.Zettel:
		if len(els) < 3 {
			err = errors.Errorf("not a valid zettel: %q, %q", els, path)
			return
		}

		if err = k2.SetWithGenre(els[2]+"/"+els[1], genres.Zettel); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		err = ErrFDNotId
		return
	}

	return
}

func (h *ObjectId) SetWithIdLike(
	k IdLike,
) (err error) {
	switch kt := k.(type) {
	case *ObjectId:
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

func (h *ObjectId) SetWithGenre(
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

func (h *ObjectId) TodoSetBytes(v *catgut.String) (err error) {
	return h.Set(v.String())
}

func (h *ObjectId) SetRaw(v string) (err error) {
	h.g = genres.Unknown

	if err = h.left.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (h *ObjectId) Set(v string) (err error) {
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

func (a *ObjectId) ResetWith(b *ObjectId) {
	a.g = b.g
	b.left.CopyTo(&a.left)
	b.right.CopyTo(&a.right)
	a.middle = b.middle
}

func (a *ObjectId) ResetWithIdLike(b IdLike) (err error) {
	return a.SetWithIdLike(b)
}

func (t *ObjectId) MarshalText() (text []byte, err error) {
	text = []byte(FormattedString(t))
	return
}

func (t *ObjectId) UnmarshalText(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (t *ObjectId) MarshalBinary() (text []byte, err error) {
	text = []byte(FormattedString(t))
	return
}

func (t *ObjectId) UnmarshalBinary(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
