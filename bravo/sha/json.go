package sha

func (s Sha) MarshalText() (text []byte, err error) {
	text = []byte(s.String())

	return
}

func (s *Sha) UnmarshalText(text []byte) (err error) {
	if err = s.Set(string(text)); err != nil {
		return
	}

	return
}
