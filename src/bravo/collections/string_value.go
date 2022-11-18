package collections

type StringValue string

func (sv *StringValue) Set(v string) (err error) {
	*sv = StringValue(v)
	return
}

func (sv StringValue) String() string {
	return string(sv)
}

func (sv StringValue) Len() int {
	return len(string(sv))
}
