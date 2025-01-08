package ids

import (
	"fmt"
	"strings"
	"unicode"

	"code.linenisgreat.com/zit/go/zit/src/alfa/coordinates"
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
)

func init() {
	register(ZettelId{})
}

type ZettelId struct {
	left, right string
}

type Provider interface {
	MakeZettelIdFromCoordinates(i coordinates.Int) (string, error)
}

func NewZettelIdEmpty() (h ZettelId) {
	h = ZettelId{}

	return
}

// TODO-P3 is this really necessary?;w
func MakeZettelIdFromProvidersAndCoordinates(
	i coordinates.Int,
	pl Provider,
	pr Provider,
) (h *ZettelId, err error) {
	k := coordinates.ZettelIdCoordinate{}
	k.SetInt(i)

	var l, r string

	if l, err = pl.MakeZettelIdFromCoordinates(k.Left); err != nil {
		err = errors.Errorf("failed to make left zettel id: %s", err)
		return
	}

	if r, err = pr.MakeZettelIdFromCoordinates(k.Right); err != nil {
		err = errors.Errorf("failed to make right zettel id: %s", err)
		return
	}

	return MakeZettelIdFromHeadAndTail(l, r)
}

func MakeZettelIdFromHeadAndTail(head, tail string) (h *ZettelId, err error) {
	head = strings.TrimSpace(head)
	tail = strings.TrimSpace(tail)

	switch {
	case head == "":
		err = errors.Errorf(
			"kopf was empty: {Kopf: %q, Schwanz: %q}",
			head,
			tail,
		)
		return

	case tail == "":
		err = errors.Errorf(
			"schwanz was empty: {Kopf: %q, Schwanz: %q}",
			head,
			tail,
		)
		return
	}

	hs := fmt.Sprintf("%s/%s", head, tail)

	h = &ZettelId{}

	if err = h.Set(hs); err != nil {
		err = errors.Errorf("failed to set zettel id: %s", err)
		return
	}

	return
}

func MustZettelId(v string) (h *ZettelId) {
	var err error
	h, err = MakeZettelId(v)

	errors.PanicIfError(err)

	return
}

func MakeZettelId(v string) (h *ZettelId, err error) {
	h = &ZettelId{}

	if err = h.Set(v); err != nil {
		return
	}

	return
}

func (a ZettelId) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a ZettelId) Equals(b ZettelId) bool {
	if a.left != b.left {
		return false
	}

	if a.right != b.right {
		return false
	}

	return true
}

func (h ZettelId) GetHead() string {
	return h.left
}

func (h ZettelId) GetTail() string {
	return h.right
}

func (i ZettelId) GetObjectIdString() string {
	return i.String()
}

func (h ZettelId) String() string {
	v := fmt.Sprintf("%s/%s", h.left, h.right)
	return v
}

func (h ZettelId) Parts() [3]string {
	return [3]string{h.left, "/", h.right}
}

func (i ZettelId) Less(j ZettelId) bool {
	return i.String() < j.String()
}

func (h *ZettelId) SetFromIdParts(parts [3]string) (err error) {
	h.left = parts[0]
	h.right = parts[2]
	return
}

func (h *ZettelId) Set(v string) (err error) {
	v = strings.TrimSpace(v)
	v = strings.ToLower(v)

	v = strings.TrimSuffix(v, ".zettel")

	me := errors.MakeMulti()

	if strings.ContainsFunc(
		v,
		func(r rune) bool {
			switch {
			case unicode.IsDigit(r), unicode.IsLetter(r), r == '_', r == '/', r == '%':
				return false

			default:
				return true
			}
		},
	) {
		me.Add(errors.Errorf("contains invalid characters: %q", v))
	}

	if v == "/" {
		if me.Len() > 0 {
			err = me
		}

		return
	}

	parts := strings.Split(v, "/")
	count := len(parts)

	switch count {
	default:
		me.Add(errors.Errorf(
			"zettel id needs exactly 2 components, but got %d: %q",
			count,
			v,
		))

	case 2:
		h.left = parts[0]
		h.right = parts[1]
	}

	if me.Len() > 0 {
		err = me
	}

	if (len(h.left) == 0 && len(h.right) > 0) ||
		(len(h.right) == 0 && len(h.left) > 0) {
		err = errors.Errorf("incomplete zettel id: %s", h)
		return
	}

	return
}

func (h *ZettelId) Reset() {
	h.left = ""
	h.right = ""
}

func (h *ZettelId) ResetWith(h1 ZettelId) {
	h.left = h1.left
	h.right = h1.right
}

func (h ZettelId) GetGenre() interfaces.Genre {
	return genres.Zettel
}

func (t ZettelId) MarshalText() (text []byte, err error) {
	text = []byte(t.String())
	return
}

func (t *ZettelId) UnmarshalText(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (t ZettelId) MarshalBinary() (text []byte, err error) {
	text = []byte(t.String())
	return
}

func (t *ZettelId) UnmarshalBinary(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
