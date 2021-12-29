package hinweis

func (s hinweis) MarshalText() (text []byte, err error) {
	text = []byte(s.String())

	return
}

func (s *hinweis) UnmarshalText(text []byte) (err error) {
	if err = s.Set(string(text)); err != nil {
		return
	}

	return
}
