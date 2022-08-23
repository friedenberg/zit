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

func (u Umwelt) FileVerzeichnisseZettelenSchwanzen() string {
	return u.DirVerzeichnisse("ZettelenSchwanzen")
}

func (u Umwelt) FileVerzeichnisseZettelen() string {
	return u.DirVerzeichnisse("Zettelen")
}

func (u Umwelt) FileVerzeichnisseEtiketten() string {
	return u.DirVerzeichnisse("Etiketten")
}

func (u Umwelt) DirVerzeichnisseAkten() string {
	return u.DirVerzeichnisse("Akten")
}
