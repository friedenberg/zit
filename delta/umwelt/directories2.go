package umwelt

func (u Umwelt) DirVerzeichnisse(p ...string) string {
	return u.DirZit(append([]string{"Verzeichnisse"}, p...)...)
}

func (u Umwelt) DirObjekten(p ...string) string {
	return u.DirZit(append([]string{"Objekten"}, p...)...)
}

func (u Umwelt) DirObjektenZettelen() string {
	return u.DirObjekten("Zettelen")
}

func (u Umwelt) DirObjektenTransaktion() string {
	return u.DirObjekten("Transaktion")
}

func (u Umwelt) DirObjektenAkten() string {
	return u.DirObjekten("Akten")
}

func (u Umwelt) DirVerzeichnisseZettelen() string {
	return u.DirVerzeichnisse("Zettelen")
}

func (u Umwelt) DirVerzeichnisseAkten() string {
	return u.DirVerzeichnisse("Akten")
}
