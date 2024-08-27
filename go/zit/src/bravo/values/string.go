package values

import "code.linenisgreat.com/zit/go/zit/src/alfa/errors"

type String struct {
	wasSet bool
	string
}

func MakeString(v string) String {
	return String{
		wasSet: true,
		string: v,
	}
}

func (sv *String) Set(v string) (err error) {
	*sv = String{
		wasSet: true,
		string: v,
	}

	return
}

func (sv String) Match(v string) (err error) {
	if sv.string != v {
		err = errors.BadRequestf("expected %q but got %q", sv.string, v)
		return
	}

	return
}

func (sv String) String() string {
	return sv.string
}

func (sv String) IsEmpty() bool {
	return len(sv.string) == 0
}

func (sv String) Len() int {
	return len(sv.string)
}

func (a String) Less(b String) bool {
	return a.string < b.string
}

func (a String) WasSet() bool {
	return a.wasSet
}

func (a *String) Reset() {
	a.wasSet = false
	a.string = ""
}

func (a *String) ResetWith(b String) {
	a.wasSet = true
	a.string = b.string
}

func (s String) MarshalBinary() (text []byte, err error) {
	text = []byte(s.String())

	return
}

func (s *String) UnmarshalBinary(text []byte) (err error) {
	if err = s.Set(string(text)); err != nil {
		return
	}

	return
}
