package bezeichnung

import "strings"

type Bezeichnung struct {
	wasSet bool
	value  string
}

func Make(v string) Bezeichnung {
	return Bezeichnung{
		wasSet: true,
		value:  v,
	}
}

func (b Bezeichnung) String() string {
	return b.value
}

func (b *Bezeichnung) Set(v string) (err error) {
	b.wasSet = true

	v1 := strings.TrimSpace(v)

	if v0 := b.String(); v0 != "" {
		b.value = v0 + " " + v1
	} else {
		b.value = v1
	}

	return
}

func (a Bezeichnung) WasSet() bool {
	return a.wasSet
}

func (a *Bezeichnung) Reset() {
	a.value = ""
}

func (a Bezeichnung) IsEmpty() bool {
	return a.value == ""
}

func (a Bezeichnung) Equals(b Bezeichnung) (ok bool) {
	if !a.wasSet {
		return false
	}

	return a.value == b.value
}

func (a Bezeichnung) Less(b Bezeichnung) (ok bool) {
	return a.value < b.value
}

func (a Bezeichnung) MarshalBinary() (text []byte, err error) {
	text = []byte(a.value)
	return
}

func (a *Bezeichnung) UnmarshalBinary(text []byte) (err error) {
	a.wasSet = true
	a.value = string(text)
	return
}

func (a Bezeichnung) MarshalText() (text []byte, err error) {
	text = []byte(a.value)
	return
}

func (a *Bezeichnung) UnmarshalText(text []byte) (err error) {
	a.wasSet = true
	a.value = string(text)
	return
}
