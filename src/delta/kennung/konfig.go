package kennung

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/values"
)

type Konfig struct{}

func (a Konfig) GetGattung() schnittstellen.Gattung {
	return gattung.Konfig
}

func (a Konfig) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a Konfig) Equals(b Konfig) bool {
	return true
}

func (a *Konfig) Reset() {
	return
}

func (a *Konfig) ResetWith(_ Konfig) {
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

func (i Konfig) GetSha() sha.Sha {
	return sha.FromString(i.String())
}
