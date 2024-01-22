package values

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type Bool struct {
	wasSet bool
	Value  bool
}

func MakeBool(v bool) Bool {
	return Bool{
		wasSet: true,
		Value:  v,
	}
}

func (sv *Bool) Set(v string) (err error) {
	v = strings.ToLower(strings.TrimSpace(v))
	v1 := false

	switch v {
	case "", "t", "true", "y", "yes":
		v1 = true
	}

	sv.SetBool(v1)

	return
}

func (sv *Bool) SetBool(v bool) {
	sv.wasSet = true
	sv.Value = v
}

func (sv Bool) Bool() bool {
	return sv.Value
}

func (sv Bool) String() string {
	return fmt.Sprintf("%t", sv.Value)
}

func (a Bool) Equals(b Bool) bool {
	return a.Value == b.Value && a.wasSet && b.wasSet
}

func (a Bool) WasSet() bool {
	return a.wasSet
}

func (a *Bool) Reset() {
	a.Value = false
	a.wasSet = false
}

func (a *Bool) ResetWith(b Bool) {
	a.wasSet = true
	a.Value = b.Value
}

func (a *Bool) MarshalBinary() ([]byte, error) {
	b := uint8(0)

	if a.Value {
		b = 1
	}

	return []byte{b}, nil
}

func (a *Bool) UnmarshalBinary(b []byte) (err error) {
	if len(b) != 1 {
		err = errors.Errorf("expected exactly 1 byte but got %d", b)
		return
	}

	b1 := b[0]

	switch b1 {
	case 0:
		a.SetBool(false)

	case 1:
		a.SetBool(true)

	default:
		err = errors.Errorf("unexpected value: %d", b1)
		return
	}

	return
}
