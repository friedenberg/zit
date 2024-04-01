package kennung

import (
	"fmt"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/coordinates"
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/values"
	"code.linenisgreat.com/zit/src/delta/gattung"
)

func init() {
	register(Hinweis{})
}

type Hinweis struct {
	left, right string
}

type Provider interface {
	Hinweis(i coordinates.Int) (string, error)
}

func NewHinweisEmpty() (h Hinweis) {
	h = Hinweis{}

	return
}

// TODO-P3 is this really necessary?;w
func NewHinweis(
	i coordinates.Int,
	pl Provider,
	pr Provider,
) (h *Hinweis, err error) {
	k := coordinates.Kennung{}
	k.SetInt(i)

	var l, r string

	if l, err = pl.Hinweis(k.Left); err != nil {
		err = errors.Errorf("failed to make left kennung: %s", err)
		return
	}

	if r, err = pr.Hinweis(k.Right); err != nil {
		err = errors.Errorf("failed to make right kennung: %s", err)
		return
	}

	return MakeHinweisKopfUndSchwanz(l, r)
}

func MakeHinweisKopfUndSchwanz(kopf, schwanz string) (h *Hinweis, err error) {
	kopf = strings.TrimSpace(kopf)
	schwanz = strings.TrimSpace(schwanz)

	switch {
	case kopf == "":
		err = errors.Errorf(
			"kopf was empty: {Kopf: %q, Schwanz: %q}",
			kopf,
			schwanz,
		)
		return

	case schwanz == "":
		err = errors.Errorf(
			"schwanz was empty: {Kopf: %q, Schwanz: %q}",
			kopf,
			schwanz,
		)
		return
	}

	hs := fmt.Sprintf("%s/%s", kopf, schwanz)

	h = &Hinweis{}

	if err = h.Set(hs); err != nil {
		err = errors.Errorf("failed to set hinweis: %s", err)
		return
	}

	return
}

func MustHinweis(v string) (h *Hinweis) {
	var err error
	h, err = MakeHinweis(v)

	errors.PanicIfError(err)

	return
}

func MakeHinweis(v string) (h *Hinweis, err error) {
	h = &Hinweis{}

	if err = h.Set(v); err != nil {
		return
	}

	return
}

func (a Hinweis) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a Hinweis) Equals(b Hinweis) bool {
	if a.left != b.left {
		return false
	}

	if a.right != b.right {
		return false
	}

	return true
}

func (h Hinweis) Kopf() string {
	return h.left
}

func (h Hinweis) Schwanz() string {
	return h.right
}

func (h Hinweis) String() string {
	v := fmt.Sprintf("%s/%s", h.left, h.right)
	return v
}

func (h Hinweis) Parts() [3]string {
	return [3]string{h.left, "/", h.right}
}

func (i Hinweis) Less(j Hinweis) bool {
	return i.String() < j.String()
}

func (h *Hinweis) Set(v string) (err error) {
	v = strings.TrimSpace(v)
	v = strings.ToLower(v)
	v = strings.Map(
		func(r rune) rune {
			if r > 'z' {
				return -1
			}

			return r
		},
		v,
	)

	v = strings.TrimSuffix(v, ".zettel")

	parts := strings.Split(v, "/")

	count := len(parts)

	switch count {
	default:
		err = errors.Errorf(
			"hinweis needs exactly 2 components, but got %d: %q",
			count,
			v,
		)
		return

	case 2:
		h.left = parts[0]
		h.right = parts[1]
	}

	switch {
	case h.left == "":
		err = errors.Errorf("left is empty: %q", v)
		return

	case h.right == "":
		err = errors.Errorf("right is empty: %q", v)
		return
	}

	return
}

func (h *Hinweis) Reset() {
	h.left = ""
	h.right = ""
}

func (h *Hinweis) ResetWith(h1 Hinweis) {
	h.left = h1.left
	h.right = h1.right
}

func (h Hinweis) GetGattung() schnittstellen.GattungLike {
	return gattung.Zettel
}

func (t Hinweis) MarshalText() (text []byte, err error) {
	text = []byte(t.String())
	return
}

func (t *Hinweis) UnmarshalText(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (t Hinweis) MarshalBinary() (text []byte, err error) {
	text = []byte(t.String())
	return
}

func (t *Hinweis) UnmarshalBinary(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
