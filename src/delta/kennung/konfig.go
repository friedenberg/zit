package kennung

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
)

type Konfig struct{}

func (a Konfig) Gattung() gattung.Gattung {
	return gattung.Konfig
}

func (a Konfig) Equals(b *Konfig) bool {
	if b == nil {
		return false
	}

	return true
}

func (a Konfig) Reset(b *Konfig) {
	return
}

func (i Konfig) String() string {
	return "konfig"
}

func (i Konfig) Set(v string) (err error) {
	v = strings.TrimSpace(v)
	v = strings.ToLower(v)

	if v != "konfig" {
		err = errors.Errorf("not konfig")
		return
	}

	return
}

func (i Konfig) Sha() sha.Sha {
	return sha.FromString(i.String())
}