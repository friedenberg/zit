package umwelt

import "path"

func (u Umwelt) Dir() string {
	return u.BasePath
}

func (u Umwelt) DirZit(p ...string) string {
	return path.Join(
		append(
			[]string{u.Dir(), ".zit"},
			p...,
		)...,
	)
}

func (u Umwelt) DirObjekte(p ...string) string {
	return u.DirZit(append([]string{"Objekte"}, p...)...)
}

func (u Umwelt) DirVerlorenUndGefunden() string {
	return u.DirZit("Verloren+Gefunden")
}

func (u Umwelt) DirZettelHinweis() string {
	return u.DirZit("Zettel-Hinweis")
}

func (u Umwelt) DirKennung() string {
	return u.DirZit("Kennung")
}

func (u Umwelt) DirHinweis() string {
	return u.DirZit("Hinweis")
}

func (u Umwelt) DirAkte() string {
	return u.DirObjekte("Akte")
}

func (u Umwelt) DirZettel() string {
	return u.DirObjekte("Zettel")
}

func (u Umwelt) FileAge() string {
	return u.DirZit("AgeIdentity")
}
