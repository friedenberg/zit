package konfig

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
)

type Id struct{}

func (i Id) String() string {
	return "konfig"
}

func (i Id) Set(v string) (err error) {
	v = strings.TrimSpace(v)
	v = strings.ToLower(v)

	if v != "konfig" {
		err = errors.Errorf("not konfig")
		return
	}

	return
}

func (i Id) Sha() sha.Sha {
	return sha.FromString(i.String())
}
