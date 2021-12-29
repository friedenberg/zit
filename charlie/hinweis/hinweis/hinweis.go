package hinweis

import (
	"fmt"
	"strings"
)

type Hinweis interface {
	_Id
}

type hinweis struct {
	left, right string
}

type Provider interface {
	Hinweis(i _Int) (string, error)
}

func New(i _Int, pl Provider, pr Provider) (h *hinweis, err error) {
	k := _Kennung{}
	k.SetInt(i)

	h = &hinweis{}

	if h.left, err = pl.Hinweis(k.Left); err != nil {
		return
	}

	if h.right, err = pr.Hinweis(k.Right); err != nil {
		return
	}

	if err = h.Set(h.String()); err != nil {
		return
	}

	return
}

func MakeBlindHinweis(v string) (h *hinweis, err error) {
	h = &hinweis{}

	if err = h.Set(v); err != nil {
		return
	}

	return
}

func MakeBlindHinweisParts(left, right string) (h hinweis) {
	h.left = left
	h.right = right

	return
}

func (h hinweis) Head() string {
	return h.left
}

func (h hinweis) Tail() string {
	return h.right
}

func (h hinweis) String() string {
	return fmt.Sprintf("%s/%s", h.left, h.right)
}

func (h *hinweis) Set(v string) (err error) {
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
	parts := strings.Split(strings.ToLower(v), "/")

	count := len(parts)

	if count != 2 {
		err = _Errorf("expected 2 components, but got %d: %s", count, v)
	}

	h.left = parts[0]
	h.right = parts[1]

	return
}
