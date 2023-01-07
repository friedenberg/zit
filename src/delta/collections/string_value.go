package collections

type StringValue struct {
	wasSet bool
	string
}

func MakeStringValue(v string) StringValue {
	return StringValue{
		wasSet: true,
		string: v,
	}
}

func (sv *StringValue) Set(v string) (err error) {
	*sv = StringValue{
		wasSet: true,
		string: v,
	}

	return
}

func (sv StringValue) String() string {
	return sv.string
}

func (sv StringValue) IsEmpty() bool {
	return len(sv.string) == 0
}

func (sv StringValue) Len() int {
	return len(sv.string)
}

func (a StringValue) Less(b StringValue) bool {
	return a.string < b.string
}

func (a StringValue) WasSet() bool {
	return a.wasSet
}

func (s StringValue) MarshalBinary() (text []byte, err error) {
	text = []byte(s.String())

	return
}

func (s *StringValue) UnmarshalBinary(text []byte) (err error) {
	if err = s.Set(string(text)); err != nil {
		return
	}

	return
}
