package hinweis

import (
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io"
	"strings"

	"github.com/friedenberg/zit/alfa/kennung"
	"github.com/friedenberg/zit/alfa/logz"
	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/bravo/stdprinter"
	"github.com/friedenberg/zit/charlie/sha"
)

func init() {
	gob.RegisterName("Hinweis", Hinweis{})
}

type Hinweis struct {
	Left, Right string
}

type Provider interface {
	Hinweis(i kennung.Int) (string, error)
}

func NewEmpty() (h Hinweis) {
	h = Hinweis{}

	return
}

//TODO is this really necessary?;w

func New(i kennung.Int, pl Provider, pr Provider) (h Hinweis, err error) {
	k := kennung.Kennung{}
	k.SetInt(i)

	h = Hinweis{}

	logz.Print("making kennung")

	logz.Print("making left")
	if h.Left, err = pl.Hinweis(k.Left); err != nil {
		logz.Printf("left failed: %s", err)
		err = errors.Errorf("failed to make left kennung: %s", err)
		return
	}

	logz.Print("making right")
	if h.Right, err = pr.Hinweis(k.Right); err != nil {
		err = errors.Errorf("failed to make right kennung: %s", err)
		return
	}

	logz.Print("making setting")
	if err = h.Set(h.String()); err != nil {
		err = errors.Errorf("failed to set hinweis: %s", err)
		return
	}

	logz.Print("done")

	return
}

func MakeBlindHinweis(v string) (h Hinweis, err error) {
	h = Hinweis{}

	if err = h.Set(v); err != nil {
		return
	}

	return
}

func MakeBlindHinweisParts(left, right string) (h Hinweis) {
	h.Left = left
	h.Right = right

	return
}

func (h Hinweis) Head() string {
	return h.Left
}

func (h Hinweis) Tail() string {
	return h.Right
}

func (h Hinweis) String() string {
	return fmt.Sprintf("%s/%s", h.Left, h.Right)
}

func (h *Hinweis) Set(v string) (err error) {
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
		err = errors.Normal(errors.Errorf("hinweis needs exactly 2 components, but got %d: %q", count, v))
		return
	}

	h.Left = parts[0]
	h.Right = parts[1]

	return
}

func (h Hinweis) Sha() sha.Sha {
	hash := sha256.New()
	sr := strings.NewReader(h.String())

	if _, err := io.Copy(hash, sr); err != nil {
		stdprinter.PanicIfError(err)
	}

	return sha.FromHash(hash)
}
