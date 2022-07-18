package hinweis

import (
	"fmt"
	"log"
	"strings"

	"github.com/friedenberg/zit/alfa/errors"
)

type Hinweis interface {
	_Id
	Set(string) error
	Equals(Hinweis) bool
}

type hinweis struct {
	left, right string
}

type Provider interface {
	Hinweis(i _Int) (string, error)
}

func NewEmpty() (h *hinweis) {
	h = &hinweis{}
	return
}

func New(i _Int, pl Provider, pr Provider) (h *hinweis, err error) {
	k := _Kennung{}
	k.SetInt(i)

	h = &hinweis{}

	log.Print("making kennung")

	log.Print("making left")
	if h.left, err = pl.Hinweis(k.Left); err != nil {
		log.Printf("left failed: %s", err)
		err = errors.Errorf("failed to make left kennung: %s", err)
		return
	}

	log.Print("making right")
	if h.right, err = pr.Hinweis(k.Right); err != nil {
		err = errors.Errorf("failed to make right kennung: %s", err)
		return
	}

	log.Print("making setting")
	if err = h.Set(h.String()); err != nil {
		err = errors.Errorf("failed to set hinweis: %s", err)
		return
	}

	log.Print("done")

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
		err = _ErrorNormal(_Errorf("hinweis needs exactly 2 components, but got %d: %q", count, v))
		return
	}

	h.left = parts[0]
	h.right = parts[1]

	return
}
