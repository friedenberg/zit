package repo_type

import (
	"fmt"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

type Type int

const (
	TypeUnknown = Type(iota)
	TypeReadWrite
	TypeRelay
)

func (t *Type) Set(v string) (err error) {
	switch strings.TrimSpace(strings.ToLower(v)) {
	case "relay":
		*t = TypeRelay

	case "", "read-write":
		*t = TypeReadWrite

	default:
		err = errors.Wrapf(ErrUnsupportedRepoType{}, "Value: %q", v)
		return
	}

	return
}

func (t Type) String() string {
	switch t {
	case TypeReadWrite:
		return "read-write"

	case TypeRelay:
		return "relay"

	default:
		return fmt.Sprintf("unknown-%d", t)
	}
}

func (t Type) MarshalText() (b []byte, err error) {
	b = []byte(t.String())
	return
}

func (t *Type) UnmarshalText(b []byte) (err error) {
	if err = t.Set(string(b)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
