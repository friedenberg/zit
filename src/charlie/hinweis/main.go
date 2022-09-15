package hinweis

import (
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/kennung"
	"github.com/friedenberg/zit/src/bravo/sha"
)

func init() {
	gob.RegisterName("Hinweis", Hinweis{})
}

type Hinweis struct {
	inner
	Bez string
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

	var l, r string

	if l, err = pl.Hinweis(k.Left); err != nil {
		err = errors.Errorf("failed to make Left kennung: %s", err)
		return
	}

	if r, err = pr.Hinweis(k.Right); err != nil {
		err = errors.Errorf("failed to make right kennung: %s", err)
		return
	}

	return MakeKopfUndSchwanz(l, r)
}

func MakeKopfUndSchwanz(kopf, schwanz string) (h Hinweis, err error) {
	hs := fmt.Sprintf("%s/%s", strings.TrimSpace(kopf), strings.TrimSpace(schwanz))

	if err = h.Set(hs); err != nil {
		err = errors.Errorf("failed to set hinweis: %s", err)
		return
	}

	return
}

func Make(v string) (h Hinweis, err error) {
	h = Hinweis{}

	if err = h.Set(v); err != nil {
		return
	}

	return
}

func (h Hinweis) String() string {
	if h.Bez == "" {
		return h.inner.String()
	} else {
		return fmt.Sprintf("%s/%s", h.inner.String(), h.Bez)
	}
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

	switch count {

	default:
		err = errors.Normal(errors.Errorf("hinweis needs exactly 2 components, but got %d: %q", count, v))
		return

	case 3:
		//Left/right/bez
		h.Bez = parts[2]
		fallthrough

	case 2:
		h.Left = parts[0]
		h.Right = parts[1]
	}

	return
}

func (h Hinweis) Sha() sha.Sha {
	hash := sha256.New()
	sr := strings.NewReader(h.inner.String())

	if _, err := io.Copy(hash, sr); err != nil {
		errors.PanicIfError(err)
	}

	return sha.FromHash(hash)
}
