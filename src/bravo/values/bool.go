package values

import (
	"fmt"
	"strings"
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

	*sv = Bool{
		wasSet: true,
		Value:  v1,
	}

	return
}

func (sv Bool) Bool() bool {
	return sv.Value
}

func (sv Bool) String() string {
	return fmt.Sprintf("%t", sv.Value)
}

func (a Bool) EqualsAny(b any) bool {
	return Equals(a, b)
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
