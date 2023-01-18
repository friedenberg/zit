package hinweis

func (s Hinweis) MarshalBinary() (text []byte, err error) {
	text = []byte(s.String())

	return
}

func (s *Hinweis) UnmarshalBinary(text []byte) (err error) {
	if err = s.Set(string(text)); err != nil {
		return
	}

	return
}

func (s Hinweis) MarshalText() (text []byte, err error) {
	text = []byte(s.String())

	return
}

func (s *Hinweis) UnmarshalText(text []byte) (err error) {
	if err = s.Set(string(text)); err != nil {
		return
	}

	return
}
