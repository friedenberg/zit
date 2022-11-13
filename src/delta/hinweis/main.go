package hinweis

import (
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/kennung"
	"github.com/friedenberg/zit/src/charlie/sha"
)

func init() {
	gob.RegisterName("Hinweis", Hinweis{})
}

type Hinweis struct {
	inner
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
	kopf = strings.TrimSpace(kopf)
	schwanz = strings.TrimSpace(schwanz)

	switch {
	case kopf == "":
		err = errors.Errorf("kopf was empty: {Kopf: %q, Schwanz: %q}", kopf, schwanz)
		return

	case schwanz == "":
		err = errors.Errorf("schwanz was empty: {Kopf: %q, Schwanz: %q}", kopf, schwanz)
		return
	}

	hs := fmt.Sprintf("%s/%s", kopf, schwanz)

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
	return h.inner.String()
}

func (h Hinweis) Aligned(kopf, schwanz int) string {
	parts := h.Parts()

	diffKopf := kopf - len(parts[0])
	if diffKopf > 0 {
		parts[0] = strings.Repeat(" ", diffKopf) + parts[0]
	}

	diffSchwanz := schwanz - len(parts[1])
	if diffSchwanz > 0 {
		parts[1] = parts[1] + strings.Repeat(" ", diffSchwanz)
	}

	return fmt.Sprintf("%s/%s", parts[0], parts[1])
}

func (h Hinweis) Parts() [2]string {
	return [2]string{h.Kopf(), h.Schwanz()}
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

	parts := strings.Split(v, "/")

	count := len(parts)

	switch count {
	default:
		err = errors.Errorf("hinweis needs exactly 2 components, but got %d: %q", count, v)
		return

	case 2:
		h.Left = parts[0]
		h.Right = parts[1]
	}

	switch {
	case h.Left == "":
		err = errors.Errorf("left is empty: %q", v)
		return

	case h.Right == "":
		err = errors.Errorf("right is empty: %q", v)
		return
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
