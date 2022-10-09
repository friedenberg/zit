package standort

import (
	"os"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/konfig"
)

type Standort struct {
	cwd      string
	basePath string
}

func Make(k konfig.Konfig) (s Standort, err error) {
	if s.basePath, err = k.DirZit(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.cwd, err = os.Getwd(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Standort) Cwd() string {
	return s.cwd
}

func (s Standort) Dir() string {
	return s.basePath
}

func (s Standort) DirZit(p ...string) string {
	return path.Join(
		append(
			[]string{s.Dir(), ".zit"},
			p...,
		)...,
	)
}

func (s Standort) FileAge() string {
	return s.DirZit("AgeIdentity")
}

func (s Standort) DirVerzeichnisse(p ...string) string {
	return s.DirZit(append([]string{"Verzeichnisse"}, p...)...)
}

func (s Standort) DirObjekten(p ...string) string {
	return s.DirZit(append([]string{"Objekten"}, p...)...)
}

func (s Standort) DirObjektenZettelen() string {
	return s.DirObjekten("Zettelen")
}

func (s Standort) DirObjektenTransaktion() string {
	return s.DirObjekten("Transaktion")
}

func (s Standort) DirObjektenAkten() string {
	return s.DirObjekten("Akten")
}

func (s Standort) DirVerlorenUndGefunden() string {
	return s.DirZit("Verloren+Gefunden")
}

func (s Standort) FileVerzeichnisseZettelenSchwanzen() string {
	return s.DirVerzeichnisse("ZettelenSchwanzen")
}

func (s Standort) DirVerzeichnisseZettelenNeue() string {
	return s.DirVerzeichnisse("ZettelenNeue")
}

func (s Standort) DirVerzeichnisseZettelenNeueSchwanzen() string {
	return s.DirVerzeichnisse("ZettelenNeueSchwanzen")
}

func (s Standort) DirVerzeichnisseAkten() string {
	return s.DirVerzeichnisse("Akten")
}

func (s Standort) FileVerzeichnisseZettelen() string {
	return s.DirVerzeichnisse("Zettelen")
}

func (s Standort) FileVerzeichnisseEtiketten() string {
	return s.DirVerzeichnisse("Etiketten")
}
