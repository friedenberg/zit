package umwelt

func (u Umwelt) DirObjekten(p ...string) string {
	return u.DirZit(append([]string{"Objekten"}, p...)...)
}

func (u Umwelt) DirZettelen() string {
	return u.DirObjekten("Zettelen")
}

func (u Umwelt) DirAkten() string {
	return u.DirObjekten("Akten")
}
